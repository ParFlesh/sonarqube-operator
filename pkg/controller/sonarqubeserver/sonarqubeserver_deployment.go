package sonarqubeserver

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
)

const (
	PodGracePeriod       int64  = 3600
	VolumePathData       string = "/opt/sonarqube/data"
	VolumePathLogs       string = "/opt/sonarqube/logs"
	VolumePathTemp       string = "/opt/sonarqube/temp"
	VolumePathExtensions string = "/opt/sonarqube/extensions"
)

type Component string

const (
	ComponentApplication Component = "application"
	ComponentSearch      Component = "search"
)

// Reconciles Deployment for SonarQubeServer
// Returns: Deployment, Error
// If Error is non-nil, Deployment is not in expected state
// Errors:
//   ErrorReasonResourceCreate: returned when Deployment does not exists
//   ErrorReasonResourceUpdate: returned when Deployment was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQubeServer) ReconcileDeployment(cr *sonarsourcev1alpha1.SonarQubeServer) (*appsv1.Deployment, error) {
	deployment, err := r.findDeployment(cr)
	if err != nil {
		return deployment, err
	}

	return deployment, nil
}

func (r *ReconcileSonarQubeServer) findDeployment(cr *sonarsourcev1alpha1.SonarQubeServer) (*appsv1.Deployment, error) {
	newDeployment, err := r.newDeployment(cr)
	if err != nil {
		return newDeployment, err
	}

	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newDeployment.Name, Namespace: newDeployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newDeployment)
		if err != nil {
			return newDeployment, err
		}
		return newDeployment, &utils.Error{
			Reason:  utils.ErrorReasonResourceCreate,
			Message: fmt.Sprintf("create Deployment %s", newDeployment.Name),
		}
	} else if err != nil {
		return newDeployment, err
	}

	return foundDeployment, nil
}

func (r *ReconcileSonarQubeServer) newDeployment(cr *sonarsourcev1alpha1.SonarQubeServer) (*appsv1.Deployment, error) {
	labels := r.Labels(cr)
	podLabels := r.PodLabels(cr)

	serviceAccount, secret, pvcs, service, err := r.getDeploymentDeps(cr)
	if err != nil {
		return nil, err
	}

	sqImage := r.getImage(cr)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cr.Name,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Replicas: cr.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    podLabels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "logs",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "temp",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "conf",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secret.Name,
									Optional:   &[]bool{true}[0],
								},
							},
						},
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcs[VolumeData].Name,
								},
							},
						},
						{
							Name: "extensions",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcs[VolumeExtensions].Name,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "sonarqube",
							Image: sqImage,
							Env: []corev1.EnvVar{
								{
									Name:  "SONARR_WEB_PORT",
									Value: fmt.Sprintf("%v", ApplicationWebPort),
								},
								{
									Name:  "SONARR_PATH_DATA",
									Value: VolumePathData,
								},
								{
									Name:  "SONARR_PATH_LOGS",
									Value: VolumePathLogs,
								},
								{
									Name:  "SONARR_PATH_TEMP",
									Value: VolumePathTemp,
								},
								{
									Name:  "SONARR_PATH_EXTENSIONS",
									Value: VolumePathExtensions,
								},
							},
							Resources: cr.Spec.Deployment.Resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: VolumePathData,
								},
								{
									Name:      "logs",
									MountPath: VolumePathLogs,
								},
								{
									Name:      "temp",
									MountPath: VolumePathTemp,
								},
								{
									Name:      "extensions",
									MountPath: VolumePathExtensions,
								},
								{
									Name:      "conf",
									MountPath: "/opt/sonarqube/conf/",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler:             corev1.Handler{},
								InitialDelaySeconds: 60,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: &[]int64{PodGracePeriod}[0],
					DNSPolicy:                     corev1.DNSDefault,
					ServiceAccountName:            serviceAccount.Name,
					Affinity: &corev1.Affinity{
						NodeAffinity:    cr.Spec.Deployment.NodeAffinity,
						PodAffinity:     cr.Spec.Deployment.PodAffinity,
						PodAntiAffinity: cr.Spec.Deployment.PodAntiAffinity,
					},
					NodeSelector:      cr.Spec.Deployment.NodeSelector,
					PriorityClassName: cr.Spec.Deployment.PriorityClass,
				},
			},
		},
	}

	switch cr.Spec.Cluster.Type {
	case sonarsourcev1alpha1.AIO, "":
		dep.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{
				Name:          "web",
				ContainerPort: ApplicationWebPort,
				Protocol:      corev1.ProtocolTCP,
			},
		}
		dep.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.TCPSocket = &corev1.TCPSocketAction{
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: ApplicationWebPort,
				StrVal: "",
			},
		}
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet = &corev1.HTTPGetAction{
			Path: "/api/system/status",
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: ApplicationWebPort,
				StrVal: "",
			},
		}
	case sonarsourcev1alpha1.Application:
		dep.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{
				Name:          "web",
				ContainerPort: ApplicationWebPort,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "node",
				ContainerPort: ApplicationPort,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "ce",
				ContainerPort: ApplicationCEPort,
				Protocol:      corev1.ProtocolTCP,
			},
		}
		hosts := cr.Spec.Cluster.Hosts
		if !utils.ContainsString(hosts, service.Spec.ClusterIP) {
			hosts = append(hosts, service.Spec.ClusterIP)
		}
		searchHosts := cr.Spec.Cluster.SearchHosts
		if !utils.ContainsString(searchHosts, service.Spec.ClusterIP) {
			searchHosts = append(searchHosts, service.Spec.ClusterIP)
		}

		clusteredEnv := []corev1.EnvVar{
			{
				Name:  "SONAR_CLUSTER_ENABLED",
				Value: strconv.FormatBool(cr.Spec.Cluster.Enabled),
			},
			{
				Name:  "SONAR_CLUSTER_NODE_TYPE",
				Value: string(cr.Spec.Cluster.Type),
			},
			{
				Name:  "SONAR_CLUSTER_SEARCH_HOSTS",
				Value: strings.Join(searchHosts, ","),
			},
			{
				Name:  "SONAR_CLUSTER_HOSTS",
				Value: strings.Join(hosts, ","),
			},
		}
		dep.Spec.Template.Spec.Containers[0].Env = append(dep.Spec.Template.Spec.Containers[0].Env, clusteredEnv...)
		dep.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.TCPSocket = &corev1.TCPSocketAction{
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: ApplicationWebPort,
				StrVal: "",
			},
		}
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet = &corev1.HTTPGetAction{
			Path: "/api/system/status",
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: ApplicationWebPort,
				StrVal: "",
			},
		}
	case sonarsourcev1alpha1.Search:
		dep.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{
				Name:          "search",
				ContainerPort: SearchPort,
				Protocol:      corev1.ProtocolTCP,
			},
		}
		searchHosts := cr.Spec.Cluster.SearchHosts
		if !utils.ContainsString(cr.Spec.Cluster.SearchHosts, service.Spec.ClusterIP) {
			searchHosts = append(searchHosts, service.Spec.ClusterIP)
		}

		clusteredEnv := []corev1.EnvVar{
			{
				Name:  "SONAR_CLUSTER_ENABLED",
				Value: strconv.FormatBool(cr.Spec.Cluster.Enabled),
			},
			{
				Name:  "SONAR_CLUSTER_NODE_TYPE",
				Value: string(cr.Spec.Cluster.Type),
			},
			{
				Name:  "SONAR_CLUSTER_SEARCH_HOSTS",
				Value: strings.Join(cr.Spec.Cluster.SearchHosts, ","),
			},
			{
				Name:  "SONAR_SEARCH_HOST",
				Value: service.Spec.ClusterIP,
			},
		}
		dep.Spec.Template.Spec.Containers[0].Env = append(dep.Spec.Template.Spec.Containers[0].Env, clusteredEnv...)
		dep.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.TCPSocket = &corev1.TCPSocketAction{
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: SearchPort,
				StrVal: "",
			},
		}
		dep.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet = &corev1.HTTPGetAction{
			Path: "/_cluster/health?wait_for_status=green&timeout=1s",
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: SearchPort,
				StrVal: "",
			},
		}
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}

func (r *ReconcileSonarQubeServer) getDeploymentDeps(cr *sonarsourcev1alpha1.SonarQubeServer) (*corev1.ServiceAccount, *corev1.Secret, map[Volume]*corev1.PersistentVolumeClaim, *corev1.Service, error) {

	serviceAccount, err := r.ReconcileServiceAccount(cr)
	if err != nil {
		return serviceAccount, nil, nil, nil, err
	}

	secret, err := r.ReconcileSecret(cr)
	if err != nil {
		return serviceAccount, secret, nil, nil, err
	}

	pvcs, err := r.ReconcilePVCs(cr)
	if err != nil {
		return serviceAccount, secret, pvcs, nil, err
	}

	service, err := r.ReconcileService(cr)
	if err != nil {
		return serviceAccount, secret, pvcs, service, err
	}

	return serviceAccount, secret, pvcs, service, nil
}
