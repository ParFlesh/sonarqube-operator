package sonarqubeserver

import (
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciles Service for SonarQubeServer
// Returns: Service, Error
// If Error is non-nil, Service is not in expected state
// Errors:
//   ErrorReasonResourceCreate: returned when Service does not exists
//   ErrorReasonResourceUpdate: returned when Service was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQubeServer) ReconcileServiceAccount(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.ServiceAccount, error) {
	foundServiceAccount, err := r.findServiceAccount(cr)
	if err != nil {
		return foundServiceAccount, err
	}

	return foundServiceAccount, nil
}

func (r *ReconcileSonarQubeServer) findServiceAccount(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.ServiceAccount, error) {
	newServiceAccount, err := r.newServiceAccount(cr)
	if err != nil {
		return newServiceAccount, err
	}

	foundServiceAccount := &corev1.ServiceAccount{}

	return foundServiceAccount, utils.CreateResourceIfNotFound(r.client, newServiceAccount, foundServiceAccount)
}

func (r *ReconcileSonarQubeServer) newServiceAccount(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.ServiceAccount, error) {
	labels := r.Labels(cr)

	dep := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cr.Name,
			Labels:    labels,
		},
	}

	if cr.Spec.Deployment.ServiceAccount != "" {
		dep.Name = cr.Spec.Deployment.ServiceAccount
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}
