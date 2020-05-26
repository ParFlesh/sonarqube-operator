package v1alpha1

import (
	"github.com/operator-framework/operator-sdk/pkg/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SonarQubeSpec defines the desired state of SonarQube
type SonarQubeSpec struct {

	// Number of SonarQube application nodes
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:number"
	Size int32 `json:"size"`

	// if empty operator will start latest version of selected edition
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Version"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	Version *string `json:"version,omitempty"`

	// community, developer, or enterprise (default is community)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Edition"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	Edition *string `json:"edition,omitempty"`

	// Automatically apply minor version updates
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Minor"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:checkbox,urn:alm:descriptor:com.tectonic.ui:advanced,urn:alm:descriptor:com.tectonic.ui:fieldGroup:updates"
	UpdatesMinor *bool `json:"updatesMinor,omitempty"`

	// Automatically apply major version updates
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Major"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:checkbox,urn:alm:descriptor:com.tectonic.ui:advanced,urn:alm:descriptor:com.tectonic.ui:fieldGroup:updates"
	UpdatesMajor *bool `json:"updatesMajor,omitempty"`

	// Secret with sonar configuration files (sonar.properties, wrapper.properties).
	// Don't add cluster properties to configuration files as this could cause unexpected results
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	Secret *string `json:"secret,omitempty"`

	// Shutdown SonarQube cluster
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Shutdown"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	Shutdown *bool `json:"shutdown,omitempty"`

	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	NodeConfig []ClusterNodeConfig `json:"nodeConfig,omitempty"`

	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	NodeConfigAdvanced []ClusterNodeConfigAdvanced `json:"nodeConfigAdvanced,omitempty"`

	// Service Account
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Service Account"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	ServiceAccount *string `json:"serviceAccount,omitempty"`
}

type ClusterNodeConfig struct {
	// Node type (all, application, or search)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Type"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:all,urn:alm:descriptor:com.tectonic.ui:select:application,urn:alm:descriptor:com.tectonic.ui:select:search,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfig,urn:alm:descriptor:com.tectonic.ui:advanced"
	Type string `json:"type"`

	// Storage class
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Storage Class"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:StorageClass,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfig,urn:alm:descriptor:com.tectonic.ui:advanced"
	StorageClass *string `json:"storageClass,omitempty"`

	// Size of Storage (ex 1Gi)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Storage Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfig,urn:alm:descriptor:com.tectonic.ui:advanced"
	StorageSize *string `json:"storageSize,omitempty"`
}

type ClusterNodeConfigAdvanced struct {
	// Node type (all, application, or search)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Type"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:all,urn:alm:descriptor:com.tectonic.ui:select:application,urn:alm:descriptor:com.tectonic.ui:select:search,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	Type string `json:"type"`

	// Node selector
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Selector"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:selector:Node,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeSelector *map[string]string `json:"nodeSelector,omitempty"`

	// Node Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:nodeAffinity,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`

	// Pod Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAffinity,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAffinity *corev1.PodAffinity `json:"podAffinity,omitempty"`

	// Pod AntiAffinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod AntiAffinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAntiAffinity,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAntiAffinity *corev1.PodAntiAffinity `json:"podAntiAffinity,omitempty"`

	// Priority Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PriorityClass *string `json:"priorityClass,omitempty"`

	// Resource requirements
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Resources"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:resourceRequirements,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:NodeConfigAdvanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// SonarQubeStatus defines the observed state of SonarQube
type SonarQubeStatus struct {
	// Conditions represent the latest available observations of an object's state
	// +optional
	Conditions status.Conditions `json:"conditions,omitempty"`

	// Kubernetes service that can be used to expose SonarQube
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Service string `json:"service,omitempty"`

	// Status of pods
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	Deployments DeploymentStatuses `json:"deployments,omitempty"`

	// Status of search pods
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Search Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	SearchDeployments DeploymentStatuses `json:"searchDeployments,omitempty"`

	// Hash of latest revision for tracking
	Revision string `json:"revision,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQube is the Schema for the sonarqubes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sonarqubes,scope=Namespaced
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="SonarQube Cluster"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="SonarQube,v1alpha1,\"\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Service,v1,\"\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Secret,v1,\"\""
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
