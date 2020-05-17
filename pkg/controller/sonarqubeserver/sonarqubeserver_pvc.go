package sonarqubeserver

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	newStatus := &sonarsourcev1alpha1.SonarQubeServerStatus{}
	*newStatus = cr.Status

	r.updateStatus(newStatus, cr)
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
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-%s", cr.Name, v), Namespace: cr.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newPVC)
		if err != nil {
			return newPVC, err
		}
		return newPVC, &utils.Error{
			Reason:  utils.ErrorReasonResourceCreate,
			Message: fmt.Sprintf("created pvc %s", newPVC.Name),
		}
	} else if err != nil {
		return newPVC, err
	}

	return foundPVC, nil
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
			err := r.client.Update(context.TODO(), cr)
			if err != nil {
				return nil, err
			}
			return nil, &utils.Error{
				Reason:  utils.ErrorReasonSpecUpdate,
				Message: "updated data storage size",
			}
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
			err := r.client.Update(context.TODO(), cr)
			if err != nil {
				return nil, err
			}
			return nil, &utils.Error{
				Reason:  utils.ErrorReasonSpecUpdate,
				Message: "updated extensions storage size",
			}
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