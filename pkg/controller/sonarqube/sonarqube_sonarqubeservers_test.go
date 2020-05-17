package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
)

// TestSonarQubeSonarQubeServers runs ReconcileSonarQube.ReconcileSonarQubeServers() against a
// fake client
func TestSonarQubeSonarQubeServers(t *testing.T) {
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
			Size: 1,
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

	// check dependencies
	for {
		_, err := r.ReconcileServiceAccount(sonarqube)
		if err != nil && utils.ReasonForError(err) == utils.ErrorReasonUnknown {
			t.Fatalf("reconcileServiceAccount: (%v)", err)
		} else if err == nil {
			break
		}
	}

	// Loop until no more errors or non-handled error
	for {
		_, err := r.ReconcileSonarQubeServers(sonarqube)
		if err != nil && utils.ReasonForError(err) == utils.ErrorReasonUnknown {
			t.Fatalf("reconcileServiceAccount: (%v)", err)
		} else if err == nil {
			break
		}
	}

	_, err := r.ReconcileSonarQubeServers(sonarqube)
	if err != nil {
		t.Error("reconcileSonarQubeServers: returned error even though SonarQubeServers is in expected state")
	}
	// check one of each type of server
	for _, v := range []sonarsourcev1alpha1.ServerType{sonarsourcev1alpha1.Application, sonarsourcev1alpha1.Search} {
		sonarQubeServer := &sonarsourcev1alpha1.SonarQubeServer{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s-%v", sonarqube.Name, v, 0), Namespace: sonarqube.Namespace}, sonarQubeServer)
		if err != nil && errors.IsNotFound(err) {
			t.Errorf("reconcileSonarQubeServers: %s SonarQubeServers not created", v)
		} else if err != nil {
			t.Fatalf("reconcileSonarQubeServers: (%v)", err)
		}
	}
}
