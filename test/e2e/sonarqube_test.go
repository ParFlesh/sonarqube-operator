package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/parflesh/sonarqube-operator/pkg/apis"
	operator "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestSonarQube(t *testing.T) {
	sonarqubeList := &operator.SonarQubeList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, sonarqubeList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("sonarqube-group", func(t *testing.T) {
		t.Run("Cluster", SonarQubeCluster)
		t.Run("Cluster2", SonarQubeCluster)
	})
}

func sonarqubeScaleTest(t *testing.T, f *framework.Framework, ctx *framework.Context) error {
	namespace, err := ctx.GetWatchNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create sonarqube custom resource
	exampleSonarQube := &operator.SonarQube{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-sonarqube",
			Namespace: namespace,
		},
		Spec: operator.SonarQubeSpec{},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleSonarQube, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	return nil
}

func SonarQubeCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewContext(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetWatchNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for sonarqube-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "sonarqube-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = sonarqubeScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
