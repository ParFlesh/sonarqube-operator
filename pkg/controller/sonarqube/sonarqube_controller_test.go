package sonarqube

import (
	"context"
	"fmt"
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

const (
	ReconcileErrorFormat string = "reconcile: (%v)"
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
	s.AddKnownTypes(sonarsourcev1alpha1.SchemeGroupVersion, sonarqube, &sonarsourcev1alpha1.SonarQubeServer{})
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
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}
	secret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Spec.Secret, Namespace: sonarqube.Namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: secret not created")
	} else if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeued even though spec should be invalid")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	if sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set to false")
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
		t.Errorf("condition invalid not set")
	}

	secret.Data["sonar.properties"] = append(secret.Data["sonar.properties"], "\nsonar.jdbc.url=test"...)
	err = r.client.Update(context.TODO(), secret)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	if sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
		t.Errorf("condition invalid not set to false")
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}
	serviceAccount := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), req.NamespacedName, serviceAccount)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: service account not created")
	} else if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue")
	}
	err = r.client.Get(context.TODO(), req.NamespacedName, sonarqube)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	if !sonarqube.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		t.Errorf("condition progressing not set")
	}
	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sonarqube.Name, Namespace: sonarqube.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		t.Error("reconcile: service not created")
	} else if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}

	// Check for search sonarqube servers
	for i := 0; i < 4; i++ {
		res, err = r.Reconcile(req)
		if err != nil {
			t.Fatalf(ReconcileErrorFormat, err)
		}
		// Check the result of reconciliation to make sure it has the desired state.
		if !res.Requeue {
			t.Error("reconcile did not requeue")
		}
	}

	for i := 0; i < 3; i++ {
		sonarQubeServer := &sonarsourcev1alpha1.SonarQubeServer{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s-%v", sonarqube.Name, sonarsourcev1alpha1.Search, i), Namespace: sonarqube.Namespace}, sonarQubeServer)
		if err != nil && errors.IsNotFound(err) {
			t.Errorf("reconcile: %s-%v sonarqube server not created", sonarsourcev1alpha1.Search, i)
		} else if err != nil {
			t.Fatalf(ReconcileErrorFormat, err)
		}
	}

	sonarQubeServer := &sonarsourcev1alpha1.SonarQubeServer{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s-%v", sonarqube.Name, sonarsourcev1alpha1.Application, 0), Namespace: sonarqube.Namespace}, sonarQubeServer)
	if err != nil && errors.IsNotFound(err) {
		t.Errorf("reconcile: %s sonarqube server not created", sonarsourcev1alpha1.Application)
	} else if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf(ReconcileErrorFormat, err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeued even though everything should be good")
	}
}
