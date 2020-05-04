package sonarqube

import (
	"context"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/status"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

func (r *ReconcileSonarQube) updateStatus(s *sonarsourcev1alpha1.SonarQubeStatus, cr *sonarsourcev1alpha1.SonarQube) {
	reqLogger := log.WithValues("SonarQube.Namespace", cr.Namespace, "SonarQube.Name", cr.Name)
	if !reflect.DeepEqual(*s, cr.Status) {
		cr.Status = *s
		err := r.client.Status().Update(context.TODO(), cr)
		if err != nil {
			reqLogger.Error(err, "failed to update status")
		}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr)
		if err != nil {
			reqLogger.Error(err, "failed to get updated sonarqube")
		}
		*s = cr.Status
	}
}

func (r *ReconcileSonarQube) ParseErrorForReconcileResult(cr *sonarsourcev1alpha1.SonarQube, err error) (reconcile.Result, error) {
	reqLogger := log.WithValues("SonarQube.Namespace", cr.Namespace, "SonarQube.Name", cr.Name)
	newStatus := cr.Status
	if err != nil && ReasonForError(err) != ErrorReasonUnknown {
		sqErr := err.(*Error)
		switch sqErr.Reason() {
		case ErrorReasonSpecUpdate, ErrorReasonResourceCreated:
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
			newStatus.Phase = sonarsourcev1alpha1.ConditionProgressing
			newStatus.Reason = newStatus.Conditions.GetCondition(sonarsourcev1alpha1.ConditionProgressing).Message
			r.updateStatus(&newStatus, cr)
			reqLogger.Info(sqErr.Error())
			return reconcile.Result{Requeue: true}, nil
		default:
			reqLogger.Error(sqErr, "unhandled sonarqube error")
			return reconcile.Result{}, sqErr
		}
	}
	return reconcile.Result{}, err
}

func isOwner(owner, child metav1.Object) bool {
	ownerUID := owner.GetUID()
	for _, v := range child.GetOwnerReferences() {
		if v.UID == ownerUID {
			return true
		}
	}
	return false
}

func (r *ReconcileSonarQube) getImage(cr *sonarsourcev1alpha1.SonarQube) string {
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
