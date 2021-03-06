package controller

import (
	"github.com/parflesh/sonarqube-operator/pkg/controller/sonarqubeserver"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sonarqubeserver.Add)
}
