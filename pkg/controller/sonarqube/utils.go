package sonarqube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isOwner(owner, child metav1.Object) bool {
	ownerUID := owner.GetUID()
	for _, v := range child.GetOwnerReferences() {
		if v.UID == ownerUID {
			return true
		}
	}
	return false
}
