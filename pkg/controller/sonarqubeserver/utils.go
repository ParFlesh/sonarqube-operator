package sonarqubeserver

import (
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/version"
	"strings"
)

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

	for k, v := range cr.Labels {
		labels[k] = v
	}

	labels[sonarsourcev1alpha1.ServerTypeLabel] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppName] = "SonarQubeServer"
	labels[sonarsourcev1alpha1.KubeAppInstance] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppVersion] = cr.Status.RevisionHash
	labels[sonarsourcev1alpha1.KubeAppManagedby] = fmt.Sprintf("sonarqube-operator.v%s", version.Version)
	labels[sonarsourcev1alpha1.KubeAppComponent] = string(cr.Spec.Type)

	return labels
}

func (r *ReconcileSonarQubeServer) PodLabels(cr *sonarsourcev1alpha1.SonarQubeServer) map[string]string {
	labels := r.Labels(cr)
	podLabels := make(map[string]string)
	podLabels[sonarsourcev1alpha1.ServerTypeLabel] = labels[sonarsourcev1alpha1.ServerTypeLabel]
	podLabels["Deployment"] = cr.Name

	return labels
}

const (
	DefaultImage = "sonarqube"
)
