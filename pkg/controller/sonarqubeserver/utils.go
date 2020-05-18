package sonarqubeserver

import (
	"context"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/status"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	"github.com/parflesh/sonarqube-operator/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

func (r *ReconcileSonarQubeServer) updateStatus(s *sonarsourcev1alpha1.SonarQubeServerStatus, cr *sonarsourcev1alpha1.SonarQubeServer) {
	reqLogger := log.WithValues("SonarQubeServer.Namespace", cr.Namespace, "SonarQubeServer.Name", cr.Name)
	if !reflect.DeepEqual(s, cr.Status) {
		cr.Status = *s
		err := r.client.Status().Update(context.TODO(), cr)
		if err != nil {
			reqLogger.Error(err, "failed to update status")
		}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
		if err != nil {
			reqLogger.Error(err, "failed to get updated sonarqube")
		}
		s = &cr.Status
	}
}

func (r *ReconcileSonarQubeServer) ParseErrorForReconcileResult(cr *sonarsourcev1alpha1.SonarQubeServer, err error) (reconcile.Result, error) {
	reqLogger := log.WithValues("SonarQubeServer.Namespace", cr.Namespace, "SonarQubeServer.Name", cr.Name)
	newStatus := cr.Status
	if err != nil && utils.ReasonForError(err) != utils.ErrorReasonUnknown {
		sqErr := err.(*utils.Error)
		switch sqErr.Type() {
		case utils.ErrorReasonSpecUpdate, utils.ErrorReasonResourceCreate, utils.ErrorReasonResourceUpdate, utils.ErrorReasonResourceWaiting:
			newStatus.Conditions.SetCondition(status.Condition{
				Type:    sonarsourcev1alpha1.ConditionProgressing,
				Status:  corev1.ConditionTrue,
				Reason:  sonarsourcev1alpha1.ConditionResourcesCreating,
				Message: sqErr.Error(),
			})
			if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
				newStatus.Conditions.SetCondition(status.Condition{
					Type:   sonarsourcev1alpha1.ConditionInvalid,
					Status: corev1.ConditionFalse,
				})
			}
			if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionPending) {
				newStatus.Conditions.SetCondition(status.Condition{
					Type:   sonarsourcev1alpha1.ConditionPending,
					Status: corev1.ConditionFalse,
				})
			}
			r.updateStatus(&newStatus, cr)
			reqLogger.Info(sqErr.Error())
			return reconcile.Result{Requeue: true}, nil
		case utils.ErrorReasonSpecInvalid, utils.ErrorReasonResourceInvalid:
			newStatus.Conditions.SetCondition(status.Condition{
				Type:    sonarsourcev1alpha1.ConditionInvalid,
				Status:  corev1.ConditionTrue,
				Reason:  sonarsourcev1alpha1.ConditionSpecInvalid,
				Message: sqErr.Error(),
			})
			if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
				newStatus.Conditions.SetCondition(status.Condition{
					Type:   sonarsourcev1alpha1.ConditionProgressing,
					Status: corev1.ConditionFalse,
				})
			}
			if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionPending) {
				newStatus.Conditions.SetCondition(status.Condition{
					Type:   sonarsourcev1alpha1.ConditionPending,
					Status: corev1.ConditionFalse,
				})
			}
			r.updateStatus(&newStatus, cr)
			reqLogger.Info(sqErr.Error())
			return reconcile.Result{}, nil
		default:
			reqLogger.Error(sqErr, "unhandled sonarqube error")
			return reconcile.Result{}, nil
		}
	}
	return reconcile.Result{}, err
}

func (r *ReconcileSonarQubeServer) getImage(cr *sonarsourcev1alpha1.SonarQubeServer) string {
	var sqImage string
	if cr.Spec.Image != "" {
		sqImage = cr.Spec.Image
	} else {
		sqImage = DefaultImage
	}

	if !strings.Contains(sqImage, ":") && cr.Spec.Version != "" {
		sqImage = fmt.Sprintf("%s:%s", sqImage, cr.Spec.Version)
	}
	return sqImage
}

func (r *ReconcileSonarQubeServer) Labels(cr *sonarsourcev1alpha1.SonarQubeServer) map[string]string {
	labels := make(map[string]string)

	labels[sonarsourcev1alpha1.ServerTypeLabel] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppName] = "SonarQubeServer"
	labels[sonarsourcev1alpha1.KubeAppInstance] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppVersion] = cr.Status.RevisionHash
	labels[sonarsourcev1alpha1.KubeAppManagedby] = fmt.Sprintf("sonarqube-operator.v%s", version.Version)
	labels[sonarsourcev1alpha1.KubeAppComponent] = string(cr.Spec.Type)

	for k, v := range cr.Labels {
		labels[k] = v
	}

	return labels
}

func (r *ReconcileSonarQubeServer) PodLabels(cr *sonarsourcev1alpha1.SonarQubeServer) map[string]string {
	labels := r.Labels(cr)
	podLabels := make(map[string]string)
	podLabels[sonarsourcev1alpha1.ServerTypeLabel] = labels[sonarsourcev1alpha1.ServerTypeLabel]
	podLabels["deployment"] = cr.Name

	return labels
}

const (
	DefaultImage = "sonarqube"
)
