package e2e

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/parflesh/sonarqube-operator/pkg/apis"
	operator "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
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

	dbDeployment, dbService := setupDependencies("database", namespace)
	err = f.Client.Create(goctx.TODO(), dbDeployment, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	err = f.Client.Create(goctx.TODO(), dbService, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, dbDeployment.Name, 1, retryInterval, timeout)
	if err != nil {
		return err
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

	// Wait for secret to be created and add sonar.jdbc.url
	for {
		sonarQube := &operator.SonarQube{}
		err := f.Client.Get(context.TODO(), types.NamespacedName{Name: exampleSonarQube.Name, Namespace: exampleSonarQube.Namespace}, sonarQube)
		if err != nil {
			return err
		}
		if sonarQube.Spec.Secret != "" {
			secret := &corev1.Secret{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: sonarQube.Spec.Secret, Namespace: sonarQube.Namespace}, secret)
			if err != nil && !errors.IsNotFound(err) {
				return err
			} else if err == nil {
				secret.Data["sonar.properties"] = append(secret.Data["sonar.properties"], "\nsonar.jdbc.url=jdbc:postgresql://postgresql/sonar?user=sonar&password=sonar"...)
				err := f.Client.Update(context.TODO(), secret)
				if err != nil {
					return err
				}
				break
			}
		}
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

	// Wait for search servers to startup
	for i := 0; i < 3; i++ {
		err := e2eutil.WaitForDeployment(t, f.KubeClient, namespace, fmt.Sprintf("%s-%s-%v", exampleSonarQube.Name, operator.Search, i), 1, retryInterval, timeout)
		if err != nil {
			return err
		}
	}

	// wait for example-sonarqube to reach 1 replica
	return nil
}

func setupDependencies(name, namespace string) (*appsv1.Deployment, *corev1.Service) {
	dbDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"database": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &[]int32{1}[0],
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"database": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-sonarqube-database",
					Namespace: namespace,
					Labels: map[string]string{
						"database": name,
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "sonarqube",
							Image: "postgres",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/var/lib/postgresql/data",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PGDATA",
									Value: "/var/lib/postgresql/data/pgdata",
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "sonar",
								},
								{
									Name:  "POSTGRES_USER",
									Value: "sonar",
								},
								{
									Name:  "POSTGRES_DB",
									Value: "sonar",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 5432,
											StrVal: "",
										},
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 5432,
											StrVal: "",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	dbService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"database": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "db",
					Protocol: corev1.ProtocolTCP,
					Port:     5432,
				},
			},
			Selector: map[string]string{
				"database": name,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	return dbDeployment, dbService
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
