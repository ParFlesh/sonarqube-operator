package sonarqubeserver

import (
	"context"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
)

// TestSonarQubeServerDeployment runs ReconcileSonarQubeServer.ReconcileDeployment() against a
// fake client
func TestSonarQubeServerPVC(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name      = "sonarqube-operator"
		namespace = "sonarqube"
	)

	// A SonarQubeServer resource with metadata and spec.
	sonarqube := &sonarsourcev1alpha1.SonarQubeServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sonarsourcev1alpha1.SonarQubeServerSpec{},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		sonarqube,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sonarsourcev1alpha1.SchemeGroupVersion, sonarqube)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClientWithScheme(s, objs...)
	// Create a ReconcileSonarQubeServer object with the scheme and fake client.
	r := &ReconcileSonarQubeServer{client: cl, scheme: s}

	_, err := r.ReconcilePVC(sonarqube)
	if utils.ReasonForError(err) != utils.ErrorReasonResourceCreate {
		t.Error("reconcileDeployment: resource created error not thrown when creating Deployment")
	}
	dataPVC := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Name, Namespace: namespace}, dataPVC)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: storage pvc not created")
	} else if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
}
