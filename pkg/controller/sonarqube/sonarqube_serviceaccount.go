package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciles Service for SonarQube
// Returns: Service, Error
// If Error is non-nil, Service is not in expected state
// Errors:
//   ErrorReasonResourceCreated: returned when Service does not exists
//   ErrorReasonResourceUpdate: returned when Service was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQube) ReconcileServiceAccount(cr *sonarsourcev1alpha1.SonarQube) (*corev1.ServiceAccount, error) {
	foundService, err := r.findServiceAccount(cr)
	if err != nil {
		return foundService, err
	}

	return foundService, nil
}

func (r *ReconcileSonarQube) findServiceAccount(cr *sonarsourcev1alpha1.SonarQube) (*corev1.ServiceAccount, error) {
	newService, err := r.newServiceAccount(cr)
	if err != nil {
		return newService, err
	}

	foundServiceAccount := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, foundServiceAccount)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newService)
		if err != nil {
			return newService, err
		}
		return newService, &Error{
			reason:  ErrorReasonResourceCreated,
			message: fmt.Sprintf("created Service %s", newService.Name),
		}
	} else if err != nil {
		return newService, err
	}

	return foundServiceAccount, nil
}

func (r *ReconcileSonarQube) newServiceAccount(cr *sonarsourcev1alpha1.SonarQube) (*corev1.ServiceAccount, error) {
	labels := r.Labels(cr)
	labels["app.kubernetes.io/component"] = "application"

	dep := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cr.Name,
			Labels:    labels,
		},
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}
