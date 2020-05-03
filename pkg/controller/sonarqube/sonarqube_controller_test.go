package sonarqube

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"testing"

	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// TestSonarQubeController runs ReconcileSonarQube.Reconcile() against a
// fake client that tracks a SonarQube object.
func TestSonarQubeController(t *testing.T) {
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

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}
	secret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Spec.Secret, Namespace: sonarqube.Namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: secret not created")
	} else if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: secret not created")
	} else if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeued even though everything should be good")
	}
}
