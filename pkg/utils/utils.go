package utils

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
