package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
)

// TestSonarQubeStatefulSet runs ReconcileSonarQube.ReconcileAppStatefulSet() against a
// fake client
func TestSonarQubeAppStatefulSet(t *testing.T) {
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
		Spec: sonarsourcev1alpha1.SonarQubeSpec{
			Secret: "test",
		},
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

	for i := 0; i < 3; i++ {
		_, err := r.ReconcileAppStatefulSet(sonarqube)
		if ReasonForError(err) != ErrorReasonResourceCreate {
			t.Error("reconcileStatefulSet: resource created error not thrown when creating StatefulSet")
		}
	}

	for i := 0; i < 3; i++ {
		_, err := r.ReconcileAppStatefulSet(sonarqube)
		if ReasonForError(err) != ErrorReasonSpecUpdate {
			t.Error("reconcileStatefulSet: resource created error not thrown when creating StatefulSet")
		}
	}

	_, err := r.ReconcileAppStatefulSet(sonarqube)
	if ReasonForError(err) != ErrorReasonResourceCreate {
		t.Error("reconcileStatefulSet: resource created error not thrown when creating StatefulSet")
	}
	statefulSet := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-app", sonarqube.Name), Namespace: sonarqube.Namespace}, statefulSet)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcileStatefulSet: StatefulSet not created")
	} else if err != nil {
		t.Fatalf("reconcileStatefulSet: (%v)", err)
	}

	statefulSet, err = r.ReconcileAppStatefulSet(sonarqube)
	if err != nil {
		t.Error("reconcileStatefulSet: returned error even though StatefulSet is in expected state")
	}
}
