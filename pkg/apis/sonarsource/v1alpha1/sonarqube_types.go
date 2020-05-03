package v1alpha1

import (
	"github.com/operator-framework/operator-sdk/pkg/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SonarQubeSpec defines the desired state of SonarQube
type SonarQubeSpec struct {
	// Version of SonarQube image to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Version"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Version string `json:"version,omitempty"`

	// Image of SonarQube to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	Image string `json:"image,omitempty"`

	// Storage
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Storage Storage `json:"storage,omitempty"`

	// Secret with sonar configuration each key will be added as environment variables into Sonar Container.
	// Refrain from adding (SONAR_CLUSTER_ENABLED, SONAR_CLUSTER_HOSTS, SONAR_CLUSTER_SEARCH_HOSTS,
	// SONAR_CLUSTER_NODE_TYPE, SONAR_AUTH_JWTBASE64HS256SECRET, SONAR_SEARCH_HOST) as these are controlled by the operator.
	// (More Information: https://docs.sonarqube.org/latest/setup/environment-variables/)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	Secret string `json:"secret,omitempty"`
}

type Storage struct {
	// Data Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Data Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Data string `json:"data,omitempty"`

	// Extensions Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Extensions Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Extensions string `json:"extensions,omitempty"`

	// Storage Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Storage Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:io.kubernetes:StorageClass"
	Class string `json:"class,omitempty"`

	// Changes volume to Empty Dir causing storage to be empty at every restart
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Ephemeral"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:checkbox"
	Ephemeral bool `json:"ephemeral,omitempty"`
}

// SonarQubeStatus defines the observed state of SonarQube
type SonarQubeStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions status.Conditions `json:"conditions,omitempty"`

	// Status of pods
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Pods PodStatuses `json:"pods,omitempty"`

	// Pod Count
	Size int32 `json:"size,omitempty"`

	// Status of instance
	Phase status.ConditionType `json:"phase,omitempty"`

	// Reason for status
	Reason string `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQube is the Schema for the sonarqubes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sonarqubes,scope=Namespaced
type SonarQube struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SonarQubeSpec   `json:"spec,omitempty"`
	Status SonarQubeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQubeList contains a list of SonarQube
type SonarQubeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SonarQube `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SonarQube{}, &SonarQubeList{})
}

type PodStatuses map[corev1.PodPhase][]string
