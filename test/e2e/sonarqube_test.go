package e2e

import (
	goctx "context"
	"fmt"
	"github.com/parflesh/sonarqube-operator/pkg/apis"
	operator "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"testing"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSonarQube(t *testing.T) {
	sonarqubeList := &operator.SonarQubeList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, sonarqubeList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("sonarqube-group", func(t *testing.T) {
		t.Run("Server", SonarQube)
	})
}

func sonarqubeDeployTest(t *testing.T, f *framework.Framework, ctx *framework.Context) error {
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
		Spec: operator.SonarQubeSpec{
			Size: 1,
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleSonarQube, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	// Wait for search servers
	for i := 0; i < 3; i++ {
		err := e2eutil.WaitForDeployment(t, f.KubeClient, namespace, fmt.Sprintf("%s-%s-%v", exampleSonarQube.Name, operator.Search, i), 0, retryInterval, timeout)
		if err != nil {
			return err
		}
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, fmt.Sprintf("%s-%s-%v", exampleSonarQube.Name, operator.Application, 0), 0, retryInterval, timeout)
	if err != nil {
		return err
	}

	// wait for example-sonarqube to reach 1 replica
	return nil
}

func SonarQube(t *testing.T) {
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

	if err = sonarqubeDeployTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
