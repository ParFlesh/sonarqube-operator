package sonarqubeserver

import (
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	ApplicationWebPort int32 = 9000
	ApplicationPort    int32 = 9003
	ApplicationCEPort  int32 = 9004
	SearchPort         int32 = 9001
)

// Reconciles Service for SonarQubeServer
// Returns: Service, Error
// If Error is non-nil, Service is not in expected state
// Errors:
//   ErrorReasonResourceCreate: returned when Service does not exists
//   ErrorReasonResourceUpdate: returned when Service was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQubeServer) ReconcileService(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.Service, error) {
	service, err := r.findService(cr)
	if err != nil {
		return service, err
	}

	newStatus := cr.Status.DeepCopy()

	newStatus.Service = service.Name

	r.updateStatus(newStatus, cr)

	if err := r.verifyService(cr, service); err != nil {
		return service, err
	}

	return service, nil
}

func (r *ReconcileSonarQubeServer) findService(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.Service, error) {
	newService, err := r.newService(cr)
	if err != nil {
		return newService, err
	}

	foundService := &corev1.Service{}

	return foundService, utils.CreateResourceIfNotFound(r.client, newService, foundService)
}

func (r *ReconcileSonarQubeServer) newService(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.Service, error) {
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
		},
	}
	switch cr.Spec.Type {
	case sonarsourcev1alpha1.AIO, "":
		dep.Spec.Ports = []corev1.ServicePort{
			{
				Name:     "web",
				Protocol: corev1.ProtocolTCP,
				Port:     ApplicationWebPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: ApplicationWebPort,
					StrVal: "",
				},
			},
		}
	case sonarsourcev1alpha1.Application:
		dep.Spec.Ports = []corev1.ServicePort{
			{
				Name:     "web",
				Protocol: corev1.ProtocolTCP,
				Port:     ApplicationWebPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: ApplicationWebPort,
					StrVal: "",
				},
			},
			{
				Name:     "ce",
				Protocol: corev1.ProtocolTCP,
				Port:     ApplicationCEPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: ApplicationCEPort,
					StrVal: "",
				},
			},
			{
				Name:     "node",
				Protocol: corev1.ProtocolTCP,
				Port:     ApplicationPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: ApplicationPort,
					StrVal: "",
				},
			},
		}
	case sonarsourcev1alpha1.Search:
		dep.Spec.Ports = []corev1.ServicePort{
			{
				Name:     "search",
				Protocol: corev1.ProtocolTCP,
				Port:     SearchPort,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: SearchPort,
					StrVal: "",
				},
			},
		}
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}

func (r *ReconcileSonarQubeServer) verifyService(cr *sonarsourcev1alpha1.SonarQubeServer, s *corev1.Service) error {
	newService, err := r.newService(cr)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(newService.Spec.Selector, s.Spec.Selector) {
		s.Spec.Selector = newService.Spec.Selector
		return utils.UpdateResource(r.client, s, utils.ErrorReasonResourceUpdate, "updated service selector")
	}

	if !reflect.DeepEqual(newService.Spec.Ports, s.Spec.Ports) {
		s.Spec.Ports = newService.Spec.Ports
		return utils.UpdateResource(r.client, s, utils.ErrorReasonResourceUpdate, "updated service ports")
	}

	if !reflect.DeepEqual(newService.Spec.Type, s.Spec.Type) {
		s.Spec.Type = newService.Spec.Type
		return utils.UpdateResource(r.client, s, utils.ErrorReasonResourceUpdate, "updated service type")
	}

	if !reflect.DeepEqual(newService.Labels, s.Labels) {
		s.Labels = newService.Labels
		return utils.UpdateResource(r.client, s, utils.ErrorReasonResourceUpdate, "updated service labels")
	}

	return nil

}
