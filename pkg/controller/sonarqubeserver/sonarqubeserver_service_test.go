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

// TestSonarQubeServerService runs ReconcileSonarQubeServer.ReconcileService() against a
// fake client
func TestSonarQubeServerService(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	var (
		namespace = "sonarqube"
	)

	// A SonarQubeServer resource with metadata and spec.
	sonarqubeList := []*sonarsourcev1alpha1.SonarQubeServer{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "server1",
				Namespace: namespace,
			},
			Spec: sonarsourcev1alpha1.SonarQubeServerSpec{},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "server2",
				Namespace: namespace,
			},
			Spec: sonarsourcev1alpha1.SonarQubeServerSpec{
				Cluster: sonarsourcev1alpha1.Cluster{
					Enabled: true,
					Type:    sonarsourcev1alpha1.Application,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "server3",
				Namespace: namespace,
			},
			Spec: sonarsourcev1alpha1.SonarQubeServerSpec{
				Cluster: sonarsourcev1alpha1.Cluster{
					Enabled: true,
					Type:    sonarsourcev1alpha1.Search,
				}},
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{}
	for _, v := range sonarqubeList {
		objs = append(objs, v)
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sonarsourcev1alpha1.SchemeGroupVersion, sonarqubeList[0])
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClientWithScheme(s, objs...)
	// Create a ReconcileSonarQubeServer object with the scheme and fake client.
	r := &ReconcileSonarQubeServer{client: cl, scheme: s}

	for _, sonarqube := range sonarqubeList {

		_, err := r.ReconcileService(sonarqube)
		if utils.ReasonForError(err) != utils.ErrorReasonResourceCreate {
			t.Error("reconcileService: resource created error not thrown when creating Service")
		}
		Service := &corev1.Service{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Name, Namespace: sonarqube.Namespace}, Service)
		if err != nil && errors.IsNotFound(err) {
			t.Error("reconcileService: Service not created")
		} else if err != nil {
			t.Fatalf("reconcileService: (%v)", err)
		}

		Service, err = r.ReconcileService(sonarqube)
		if err != nil {
			t.Error("reconcileService: returned error even though Service is in expected state")
		}
	}
}
