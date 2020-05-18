package utils

import (
	"fmt"
	"github.com/magiconair/properties"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsOwner(owner, child metav1.Object) bool {
	ownerUID := owner.GetUID()
	for _, v := range child.GetOwnerReferences() {
		if v.UID == ownerUID {
			return true
		}
	}
	return false
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetProperties(s *corev1.Secret, f string) (*properties.Properties, error) {
	if v, ok := s.Data[f]; ok {
		return properties.Load(v, properties.UTF8)
	} else {
		return nil, &Error{
			Reason:  ErrorReasonSpecInvalid,
			Message: fmt.Sprintf("%s doesn't exist in secret %s", f, s.Name),
		}
	}
}

func GetDeploymentCondition(deployment *appsv1.Deployment, condition appsv1.DeploymentConditionType) corev1.ConditionStatus {
	for _, v := range deployment.Status.Conditions {
		if v.Type == condition {
			return v.Status
		}
	}
	return corev1.ConditionUnknown
}
