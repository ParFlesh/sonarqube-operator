package sonarqubeserver

import (
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciles map[Volume]*PersistentVolumeClaim for SonarQubeServer
// Returns: map[Volume]*PersistentVolumeClaim, Error
// If Error is non-nil, map[Volume]*PersistentVolumeClaim is not in expected state
// Errors:
//   ErrorReasonResourceCreate: returned when any PersistentVolumeClaim does not exists
//   ErrorReasonResourceUpdate: returned when any PersistentVolumeClaim was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQubeServer) ReconcilePVCs(cr *sonarsourcev1alpha1.SonarQubeServer) (map[Volume]*corev1.PersistentVolumeClaim, error) {
	pvcs, err := r.findPVCs(cr)
	if err != nil {
		return pvcs, err
	}

	newStatus := cr.DeepCopy()

	utils.UpdateStatus(r.client, newStatus, cr)
	return pvcs, nil
}

func (r *ReconcileSonarQubeServer) findPVCs(cr *sonarsourcev1alpha1.SonarQubeServer) (map[Volume]*corev1.PersistentVolumeClaim, error) {
	foundPVCs := make(map[Volume]*corev1.PersistentVolumeClaim)

	for _, v := range []Volume{VolumeData, VolumeExtensions} {
		pvc, err := r.findPVC(cr, v)
		if err != nil {
			return foundPVCs, err
		}
		foundPVCs[v] = pvc
	}

	return foundPVCs, nil
}

func (r *ReconcileSonarQubeServer) findPVC(cr *sonarsourcev1alpha1.SonarQubeServer, v Volume) (*corev1.PersistentVolumeClaim, error) {
	newPVC, err := r.newPVC(cr, v)
	if err != nil {
		return newPVC, err
	}

	foundPVC := &corev1.PersistentVolumeClaim{}

	return foundPVC, utils.CreateResourceIfNotFound(r.client, newPVC, foundPVC)
}

func (r *ReconcileSonarQubeServer) newPVC(cr *sonarsourcev1alpha1.SonarQubeServer, v Volume) (*corev1.PersistentVolumeClaim, error) {
	labels := r.Labels(cr)

	dep := &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", cr.Name, v),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{},
			},
			VolumeMode: &[]corev1.PersistentVolumeMode{corev1.PersistentVolumeFilesystem}[0],
		},
	}

	switch v {
	case VolumeData:
		if cr.Spec.Storage.DataSize == "" {
			cr.Spec.Storage.DataSize = DefaultVolumeSize
			return nil, utils.UpdateResource(r.client, cr, utils.ErrorReasonSpecUpdate, "updated data storage size")
		}
		dep.Spec.StorageClassName = cr.Spec.Storage.DataClass
		if size, err := resource.ParseQuantity(cr.Spec.Storage.DataSize); err != nil {
			return nil, err
		} else {
			dep.Spec.Resources.Requests[corev1.ResourceStorage] = size
		}
	case VolumeExtensions:
		if cr.Spec.Storage.ExtensionsSize == "" {
			cr.Spec.Storage.ExtensionsSize = DefaultVolumeSize
			return nil, utils.UpdateResource(r.client, cr, utils.ErrorReasonSpecUpdate, "updated extensions storage size")
		}
		dep.Spec.StorageClassName = cr.Spec.Storage.ExtensionsClass
		if size, err := resource.ParseQuantity(cr.Spec.Storage.ExtensionsSize); err != nil {
			return nil, err
		} else {
			dep.Spec.Resources.Requests[corev1.ResourceStorage] = size
		}
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}

type Volume string

const (
	VolumeData        Volume = "data"
	VolumeExtensions  Volume = "extensions"
	DefaultVolumeSize        = "1Gi"
)
