package sonarqubeserver

import (
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
)

func (r *ReconcileSonarQubeServer) ReconcileServer(cr *sonarsourcev1alpha1.SonarQubeServer) error {
	service, err := r.ReconcileService(cr)
	if err != nil {
		return err
	}
	apiClient := r.apiClient.New(fmt.Sprintf("http://%s:%v", service.Spec.ClusterIP, service.Spec.Ports[0].Port))

	err = apiClient.Ping()
	if err != nil {
		return &utils.Error{
			Reason:  utils.ErrorReasonServerWaiting,
			Message: fmt.Sprintf("waiting for api to respond (%s)", err.Error()),
		}
	}
	return nil
}
