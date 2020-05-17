package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciles Service for SonarQube
// Returns: Service, Error
// If Error is non-nil, Service is not in expected state
// Errors:
//   ErrorReasonResourceCreate: returned when Service does not exists
//   ErrorReasonResourceUpdate: returned when Service was updated to meet expected state
//   ErrorReasonUnknown: returned when unhandled error from client occurs
func (r *ReconcileSonarQube) ReconcileSonarQubeServers(cr *sonarsourcev1alpha1.SonarQube) (map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer, error) {
	sonarQubeServers, err := r.findSonarQubeServers(cr)
	if err != nil {
		return sonarQubeServers, err
	}

	return sonarQubeServers, nil
}

func (r *ReconcileSonarQube) findSonarQubeServers(cr *sonarsourcev1alpha1.SonarQube) (map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer, error) {
	sonarQubeServers := make(map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer)
	newSonarQubeServers, err := r.newSonarQubeServers(cr)
	if err != nil {
		return sonarQubeServers, err
	}

	for t, l := range newSonarQubeServers {
		for _, s := range l {
			sonarQubeServer := &sonarsourcev1alpha1.SonarQubeServer{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: s.Name, Namespace: cr.Namespace}, sonarQubeServer)
			if err != nil && errors.IsNotFound(err) {
				err := r.client.Create(context.TODO(), s)
				if err != nil {
					return sonarQubeServers, err
				}
				sonarQubeServers[t] = append(sonarQubeServers[t], sonarQubeServer)
				return sonarQubeServers, &utils.Error{
					Reason:  utils.ErrorReasonResourceCreate,
					Message: fmt.Sprintf("created sonarqube server %s", s.Name),
				}
			} else if err != nil {
				return sonarQubeServers, err
			}
			sonarQubeServers[t] = append(sonarQubeServers[t], sonarQubeServer)
		}
	}

	return sonarQubeServers, nil
}

func (r *ReconcileSonarQube) newSonarQubeServers(cr *sonarsourcev1alpha1.SonarQube) (map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer, error) {
	sonarQubeServers := make(map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer)

	var i int32
	for i = 0; i < 3; i++ {
		dep, err := r.newSonarQubeServer(cr, sonarsourcev1alpha1.Search, i)
		if err != nil {
			return sonarQubeServers, err
		}
		sonarQubeServers[sonarsourcev1alpha1.Search] = append(sonarQubeServers[sonarsourcev1alpha1.Search], dep)
	}

	for i = 0; i < cr.Spec.Size; i++ {
		dep, err := r.newSonarQubeServer(cr, sonarsourcev1alpha1.Application, i)
		if err != nil {
			return sonarQubeServers, err
		}
		sonarQubeServers[sonarsourcev1alpha1.Application] = append(sonarQubeServers[sonarsourcev1alpha1.Application], dep)
	}

	return sonarQubeServers, nil
}

func (r *ReconcileSonarQube) newSonarQubeServer(cr *sonarsourcev1alpha1.SonarQube, component sonarsourcev1alpha1.ServerType, i int32) (*sonarsourcev1alpha1.SonarQubeServer, error) {
	labels := r.Labels(cr)
	labels[sonarsourcev1alpha1.KubeAppComponent] = string(component)

	serviceAccount, err := r.ReconcileServiceAccount(cr)
	if err != nil {
		return nil, err
	}

	dep := &sonarsourcev1alpha1.SonarQubeServer{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%v", cr.Name, component, i),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: sonarsourcev1alpha1.SonarQubeServerSpec{
			Size:    &[]int32{0}[0],
			Version: cr.Spec.Version,
			Image:   cr.Spec.Image,
			Secret:  cr.Spec.Secret,
			Cluster: sonarsourcev1alpha1.Cluster{
				Enabled:     true,
				Type:        component,
				Hosts:       nil,
				SearchHosts: nil,
			},
			Deployment: sonarsourcev1alpha1.Deployment{
				ServiceAccount: serviceAccount.Name,
			},
			Storage: sonarsourcev1alpha1.Storage{},
		},
	}

	if err := controllerutil.SetControllerReference(cr, dep, r.scheme); err != nil {
		return dep, err
	}

	return dep, nil
}
