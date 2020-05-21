package sonarqube

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	"github.com/thanhpk/randstr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

// Reconciles Secret for SonarQube
// Returns: Secret, Error
// If Error is non-nil, Service is not in expected state
// Errors:
//   ErrorReasonSpecUpdate: returned when spec does not have secret name
//   ErrorReasonResourceCreate: returned when secret does not exists
//   ErrorReasonResourceUpdate: returned when secret was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQube) ReconcileSecret(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Secret, error) {
	foundSecret, err := r.findSecret(cr)
	if err != nil {
		return foundSecret, err
	}

	if !isOwner(cr, foundSecret) {
		annotations := foundSecret.GetAnnotations()
		if val, ok := annotations[sonarsourcev1alpha1.SecretAnnotation]; ok && !strings.Contains(val, cr.Name) {
			annotations[sonarsourcev1alpha1.SecretAnnotation] = fmt.Sprintf("%s,%s", val, cr.Name)
			foundSecret.SetAnnotations(annotations)
			return foundSecret, utils.UpdateResource(r.client, foundSecret, utils.ErrorReasonResourceUpdate, "updated secret annotation")
		} else if !ok {
			if annotations == nil {
				annotations = make(map[string]string)
			}
			annotations[sonarsourcev1alpha1.SecretAnnotation] = cr.Name
			foundSecret.SetAnnotations(annotations)
			return foundSecret, utils.UpdateResource(r.client, foundSecret, utils.ErrorReasonResourceUpdate, "updated secret annotation")
		}
	}

	err = r.verifySecret(cr, foundSecret)
	if err != nil {
		return foundSecret, err
	}

	return foundSecret, nil
}

func (r *ReconcileSonarQube) findSecret(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Secret, error) {
	newSecret, err := r.newSecret(cr)
	if err != nil {
		return newSecret, err
	}

	foundSecret := &corev1.Secret{}

	return foundSecret, utils.CreateResourceIfNotFound(r.client, newSecret, foundSecret)
}

func (r *ReconcileSonarQube) newSecret(cr *sonarsourcev1alpha1.SonarQube) (*corev1.Secret, error) {
	labels := r.Labels(cr)

	dep := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"sonar.properties": []byte(""),
			"wrapper.conf":     []byte(""),
		},
		Type: corev1.SecretTypeOpaque,
	}

	if cr.Spec.Secret == "" {
		cr.Spec.Secret = fmt.Sprintf("%s-config", cr.Name)
		return dep, utils.UpdateResource(r.client, cr, utils.ErrorReasonSpecUpdate, "updated secret")
	}

	dep.Name = cr.Spec.Secret

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}

func (r *ReconcileSonarQube) verifySecret(cr *sonarsourcev1alpha1.SonarQube, s *corev1.Secret) error {
	sonarProperties, err := utils.GetProperties(s, "sonar.properties")
	if err != nil {
		return err
	}

	if _, ok := sonarProperties.Get("sonar.jdbc.url"); !ok {
		return &utils.Error{
			Reason:  utils.ErrorReasonSpecInvalid,
			Message: "sonar.jdbc.url not set",
		}
	}

	if _, ok := sonarProperties.Get("sonar.auth.jwtBase64Hs256Secret"); !ok && isOwner(cr, s) {
		secret := randstr.String(8)
		data := randstr.String(32)

		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, []byte(secret))

		// Write Data to it
		h.Write([]byte(data))

		// Get result and encode as hexadecimal string
		sha := hex.EncodeToString(h.Sum(nil))

		s.Data["sonar.properties"] = append(s.Data["sonar.properties"], "\nsonar.auth.jwtBase64Hs256Secret="...)
		s.Data["sonar.properties"] = append(s.Data["sonar.properties"], sha...)

		return utils.UpdateResource(r.client, s, utils.ErrorReasonResourceUpdate, fmt.Sprintf("added sonar.auth.jwtBase64Hs256Secret to sonar.properties in %s", s.Name))
	} else if !ok {
		// Don't make changes to unowned resources
		return &utils.Error{
			Reason:  utils.ErrorReasonSpecInvalid,
			Message: "sonar.auth.jwtBase64Hs256Secret not set",
		}
	}

	return nil
}
