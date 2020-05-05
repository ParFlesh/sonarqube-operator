package sonarqube

import (
	"context"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
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

// TestSonarQubeService runs ReconcileSonarQube.ReconcileAppService() against a
// fake client
func TestSonarQubeService(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		name      = "sonarqube-operator"
		namespace = "sonarqube"
	)

	// A SonarQube resource with metadata and spec.
	sonarqube := &sonarsourcev1alpha1.SonarQube{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sonarsourcev1alpha1.SonarQubeSpec{},
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
	// Create a ReconcileSonarQube object with the scheme and fake client.
	r := &ReconcileSonarQube{client: cl, scheme: s}

	_, err := r.ReconcileAppService(sonarqube)
	if ReasonForError(err) != ErrorReasonResourceCreate {
		t.Error("reconcileService: resource created error not thrown when creating Service")
	}
	Service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Name, Namespace: sonarqube.Namespace}, Service)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcileService: Service not created")
	} else if err != nil {
		t.Fatalf("reconcileService: (%v)", err)
	}

	Service, err = r.ReconcileAppService(sonarqube)
	if err != nil {
		t.Error("reconcileService: returned error even though Service is in expected state")
	}
}
