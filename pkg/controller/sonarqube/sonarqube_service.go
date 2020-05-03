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

const (
	SonarQubePort int32 = 9000
)

// Reconciles Service for SonarQube
// Returns Service if Service is in the expected state
// Errors:
// ErrorReasonResourceCreated: returned when Service does not exists
// ErrorReasonResourceUpdate: returned when Service was updated to meet expected state
// ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQube) ReconcileService(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Service, error) {
	foundService, err := r.findService(cr)
	if err != nil {
		return foundService, err
	}

	return foundService, nil
}

func (r *ReconcileSonarQube) findService(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Service, error) {
	newService, err := r.newService(cr)
	if err != nil {
		return newService, err
	}

	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newService.Name, Namespace: newService.Namespace}, foundService)
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

	return foundService, nil
}

func (r *ReconcileSonarQube) newService(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Service, error) {
	labels := r.Labels(cr)

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cr.Name,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     SonarQubePort,
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}
