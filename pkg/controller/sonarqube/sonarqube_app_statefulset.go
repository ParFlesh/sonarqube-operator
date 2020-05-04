package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciles StatefulSet for SonarQube
// Returns: StatefulSet, Error
// If Error is non-nil, StatefulSet is not in expected state
// Errors:
//   ErrorReasonResourceCreated: returned when StatefulSet does not exists
//   ErrorReasonResourceUpdate: returned when StatefulSet was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQube) ReconcileAppStatefulSet(cr *sonarsourcev1alpha1.SonarQube) (*appsv1.StatefulSet, error) {
	foundStatefulSet, err := r.findAppStatefulSet(cr)
	if err != nil {
		return foundStatefulSet, err
	}

	return foundStatefulSet, nil
}

const (
	PodGracePeriod       int64  = 3600
	VolumePathData       string = "/opt/sonarqube/data"
	VolumePathLogs       string = "/opt/sonarqube/logs"
	VolumePathTemp       string = "/opt/sonarqube/temp"
	VolumePathExtensions string = "/opt/sonarqube/extensions"
)

func (r *ReconcileSonarQube) findAppStatefulSet(cr *sonarsourcev1alpha1.SonarQube) (*appsv1.StatefulSet, error) {
	newStatefulSet, err := r.newAppStatefulSet(cr)
	if err != nil {
		return newStatefulSet, err
	}

	foundStatefulSet := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newStatefulSet.Name, Namespace: newStatefulSet.Namespace}, foundStatefulSet)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newStatefulSet)
		if err != nil {
			return newStatefulSet, err
		}
		return newStatefulSet, &Error{
			reason:  ErrorReasonResourceCreated,
			message: fmt.Sprintf("created StatefulSet %s", newStatefulSet.Name),
		}
	} else if err != nil {
		return newStatefulSet, err
	}

	return foundStatefulSet, nil
}

func (r *ReconcileSonarQube) newAppStatefulSet(cr *sonarsourcev1alpha1.SonarQube) (*appsv1.StatefulSet, error) {
	labels := r.Labels(cr)
	labels["app.kubernetes.io/component"] = "application"

	serviceAccount, service, secret, err := r.getAppStatefulSetDeps(cr)
	if err != nil {
		return &appsv1.StatefulSet{}, err
	}

	var dataVolumeRequest, extensionsVolumeRequest resource.Quantity
	if cr.Spec.Node.Storage.Data != "" {
		if dataVolumeRequest, err = resource.ParseQuantity(cr.Spec.Node.Storage.Data); err != nil {
			return &appsv1.StatefulSet{}, err
		}
		if extensionsVolumeRequest, err = resource.ParseQuantity(cr.Spec.Node.Storage.Extensions); err != nil {
			return &appsv1.StatefulSet{}, err
		}
	} else {
		if dataVolumeRequest, err = resource.ParseQuantity(DefaultVolumeSize); err != nil {
			return &appsv1.StatefulSet{}, err
		}
		if extensionsVolumeRequest, err = resource.ParseQuantity(DefaultVolumeSize); err != nil {
			return &appsv1.StatefulSet{}, err
		}
	}

	sqImage := r.getImage(cr)

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      fmt.Sprintf("%s-app", cr.Name),
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &cr.Spec.Node.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    labels,
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
					},
					Containers: []corev1.Container{
						{
							Name:  "sonarqube",
							Image: sqImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: ApplicationPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: secret.Name,
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SONARR_WEB_PORT",
									Value: string(ApplicationPort),
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
							Resources: cr.Spec.Node.Resources,
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
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: ApplicationPort,
											StrVal: "",
										},
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/api/system/status",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: ApplicationPort,
											StrVal: "",
										},
									},
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: &[]int64{PodGracePeriod}[0],
					DNSPolicy:                     corev1.DNSDefault,
					ServiceAccountName:            serviceAccount.Name,
					Affinity: &corev1.Affinity{
						NodeAffinity:    cr.Spec.Node.NodeAffinity,
						PodAffinity:     cr.Spec.Node.PodAffinity,
						PodAntiAffinity: cr.Spec.Node.PodAntiAffinity,
					},
					NodeSelector:      cr.Spec.Node.NodeSelector,
					PriorityClassName: cr.Spec.Node.PriorityClass,
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "data",
						Namespace: cr.Namespace,
						Labels:    labels,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Limits: nil,
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: dataVolumeRequest,
							},
						},
						VolumeMode: &[]corev1.PersistentVolumeMode{corev1.PersistentVolumeFilesystem}[0],
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "extensions",
						Namespace: cr.Namespace,
						Labels:    labels,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Limits: nil,
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: extensionsVolumeRequest,
							},
						},
						VolumeMode: &[]corev1.PersistentVolumeMode{corev1.PersistentVolumeFilesystem}[0],
					},
				},
			},
			ServiceName:         service.Name,
			PodManagementPolicy: appsv1.OrderedReadyPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}

	if cr.Spec.Node.Storage.Class != "" {
		dep.Spec.VolumeClaimTemplates[0].Spec.StorageClassName = &cr.Spec.Node.Storage.Class
		dep.Spec.VolumeClaimTemplates[1].Spec.StorageClassName = &cr.Spec.Node.Storage.Class
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}

func (r *ReconcileSonarQube) getAppStatefulSetDeps(cr *sonarsourcev1alpha1.SonarQube) (*corev1.ServiceAccount, *corev1.Service, *corev1.Secret, error) {

	serviceAccount, err := r.ReconcileServiceAccount(cr)
	if err != nil {
		return serviceAccount, &corev1.Service{}, &corev1.Secret{}, err
	}

	service, err := r.ReconcileAppService(cr)
	if err != nil {
		return serviceAccount, service, &corev1.Secret{}, err
	}

	secret, err := r.ReconcileSecret(cr)
	if err != nil {
		return serviceAccount, service, secret, err
	}

	return serviceAccount, service, secret, nil
}
