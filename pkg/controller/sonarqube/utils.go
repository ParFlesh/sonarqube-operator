package sonarqube

import (
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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

// getPodStatuses returns the map of pod names of the array of pods passed in
func getPodStatuses(pods []corev1.Pod) map[corev1.PodPhase][]string {
	podStatuses := make(sonarsourcev1alpha1.PodStatuses)
	for _, pod := range pods {
		if v, ok := podStatuses[pod.Status.Phase]; ok {
			podStatuses[pod.Status.Phase] = append(v, pod.Name)
		} else {
			podStatuses[pod.Status.Phase] = []string{pod.Name}
		}
	}
	return podStatuses
}
