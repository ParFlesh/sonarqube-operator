package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
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

	err = r.verifySonarQubeServers(cr, sonarQubeServers)
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
			err := utils.CreateResourceIfNotFound(r.client, s, sonarQubeServer)
			if err != nil {
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
	labels[sonarsourcev1alpha1.KubeAppPartof] = cr.Name

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
			Size:        0,
			Version:     cr.Spec.Version,
			Image:       cr.Spec.Image,
			Secret:      cr.Spec.Secret,
			Type:        component,
			Hosts:       nil,
			SearchHosts: nil,
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

func (r *ReconcileSonarQube) verifySonarQubeServers(cr *sonarsourcev1alpha1.SonarQube, s map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer) error {
	// Wait for all resources to be ready and valid to continue

	for t, l := range s {
		for _, v := range l {
			if !v.Status.Conditions.IsFalseFor(sonarsourcev1alpha1.ConditionProgressing) && !v.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
				return &utils.Error{
					Reason:  utils.ErrorReasonResourceWaiting,
					Message: fmt.Sprintf("waiting for %s node %s to finish startup", t, v.Name),
				}
			} else if v.Status.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
				return &utils.Error{
					Reason:  utils.ErrorReasonResourceInvalid,
					Message: fmt.Sprintf("%s node %s is invalid: %s", t, v.Name, v.Status.Conditions.GetCondition(sonarsourcev1alpha1.ConditionInvalid).Message),
				}
			}
		}
	}

	err := r.verifySonarQubeServersSearchHosts(cr, s)
	if err != nil {
		return err
	}

	if cr.Spec.Shutdown {
		err := r.shutdownCluster(cr, s)
		if err != nil {
			return err
		}
	} else {
		err := r.startupCluster(cr, s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileSonarQube) shutdownCluster(_ *sonarsourcev1alpha1.SonarQube, _ map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer) error {
	return nil
}

func (r *ReconcileSonarQube) startupCluster(_ *sonarsourcev1alpha1.SonarQube, s map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer) error {
	for _, v := range s[sonarsourcev1alpha1.Search] {
		if v.Spec.Size != 1 {
			v.Spec.Size = 1
			return utils.UpdateResource(r.client, v, utils.ErrorReasonResourceUpdate, fmt.Sprintf("starting sonarqube server %s", v.Name))
		}
	}

	for _, v := range s[sonarsourcev1alpha1.Application] {
		if v.Spec.Size != 1 {
			v.Spec.Size = 1
			return utils.UpdateResource(r.client, v, utils.ErrorReasonResourceUpdate, fmt.Sprintf("starting sonarqube server %s", v.Name))
		}
	}

	return nil
}

func (r *ReconcileSonarQube) getSonarQubeServersClusterIP(s []*sonarsourcev1alpha1.SonarQubeServer) ([]string, error) {
	var ips []string

	for _, v := range s {
		if v.Status.Service == "" {
			return ips, &utils.Error{
				Reason:  utils.ErrorReasonResourceWaiting,
				Message: fmt.Sprintf("Waiting on service for %s", v.Name),
			}
		}
		service := &corev1.Service{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: v.Status.Service, Namespace: v.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			return ips, &utils.Error{
				Reason:  utils.ErrorReasonResourceWaiting,
				Message: fmt.Sprintf("Waiting on service for %s", v.Name),
			}
		} else if err != nil {
			return ips, err
		}
		if service.Spec.ClusterIP == "" {
			return ips, &utils.Error{
				Reason:  utils.ErrorReasonResourceWaiting,
				Message: fmt.Sprintf("Waiting on service for %s", v.Name),
			}
		}
		ips = append(ips, service.Spec.ClusterIP)
	}

	return ips, nil
}

func (r *ReconcileSonarQube) verifySonarQubeServersSearchHosts(_ *sonarsourcev1alpha1.SonarQube, s map[sonarsourcev1alpha1.ServerType][]*sonarsourcev1alpha1.SonarQubeServer) error {
	searchServiceIPS, err := r.getSonarQubeServersClusterIP(s[sonarsourcev1alpha1.Search])
	if err != nil {
		return err
	}

	applicationServiceIPS, err := r.getSonarQubeServersClusterIP(s[sonarsourcev1alpha1.Application])
	if err != nil {
		return err
	}

	for t, l := range s {
		for _, v := range l {
			var update bool
			if !reflect.DeepEqual(v.Spec.SearchHosts, searchServiceIPS) {
				v.Spec.SearchHosts = searchServiceIPS
				update = true
			}
			if t == sonarsourcev1alpha1.Application && !reflect.DeepEqual(v.Spec.Hosts, applicationServiceIPS) {
				v.Spec.Hosts = applicationServiceIPS
				update = true
			}
			if update {
				return utils.UpdateResource(r.client, v, utils.ErrorReasonResourceUpdate, fmt.Sprintf("updates sonarqube server %s", v.Name))
			}
		}
	}

	return nil
}
