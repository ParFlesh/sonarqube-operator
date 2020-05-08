package v1alpha1

import (
	"github.com/operator-framework/operator-sdk/pkg/status"
)

// Condition Types
const (
	// ConditionPending means resources have been created, but one or more resources are not running/ready.
	ConditionPending status.ConditionType = "Pending"
	// ConditionRunning means the instance has been created and all of the resources are running/ready.
	ConditionRunning status.ConditionType = "Running"
	// ConditionInvalid means that there is a misconfiguration that can not be corrected by the operator.
	ConditionInvalid status.ConditionType = "Invalid"
	// ConditionProgressing means that for some reason the state of the resources did not match the expected state.
	// Resources are being updated to meet expected state.
	ConditionProgressing status.ConditionType = "Progressing"
)

// Condition Reasons
const (
	// ConditionResourcesCreating means that resources are being created
	ConditionResourcesCreating status.ConditionReason = "CreatingResources"
	// ConditionReasourceUpdating means that resources are updating
	ConditionReasourcesUpdating status.ConditionReason = "ResourcesUpdating"
)

const (
	SecretAnnotation = "sonarqube.sonarsource.parflesh.github.io/database"
)
