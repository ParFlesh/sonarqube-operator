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

	// Number of SonarQube Compute Engines
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:number"
	Size int32 `json:"size,omitempty"`

	// Version of SonarQube image to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Version"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Version string `json:"version,omitempty"`

	// Image of SonarQube to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	Image string `json:"image,omitempty"`

	// Secret with sonar configuration files (sonar.properties).
	// Don't add cluster properties to configuration files as this could cause unexpected results
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	Secret string `json:"secret,omitempty"`

	// Run SonarQube as a cluster
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Clustered"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Clustered bool `json:"clustered,omitempty"`

	// Shutdown SonarQube cluster
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Shutdown"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Shutdown bool `json:"shutdown,omitempty"`

	// Pod Configuration
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Node ApplicationPodConfig `json:"node,omitempty"`

	// Pod Configuration
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	NodeSearch SearchPodConfig `json:"nodeSearch,omitempty"`
}

type ApplicationPodConfig struct {
	// Node selector
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Selector"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:selector:Node,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Node Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:nodeAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`

	// Pod Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAffinity *corev1.PodAffinity `json:"podAffinity,omitempty"`

	// Pod AntiAffinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod AntiAffinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAntiAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAntiAffinity *corev1.PodAntiAffinity `json:"podAntiAffinity,omitempty"`

	// Priority Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PriorityClass string `json:"priorityClass,omitempty"`

	// Resource requirements
	// +optional
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:resourceRequirements,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Storage
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Storage ApplicationStorage `json:"storage,omitempty"`
}

type ApplicationStorage struct {
	// Data Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Data Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.storage"
	Data string `json:"data,omitempty"`

	// Extensions Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Extensions Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.storage"
	Extensions string `json:"extensions,omitempty"`

	// Storage Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Storage Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:StorageClass,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.application.storage"
	Class string `json:"class,omitempty"`
}
type SearchPodConfig struct {
	// Node selector
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Selector"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:selector:Node,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Node Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:nodeAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`

	// Pod Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAffinity *corev1.PodAffinity `json:"podAffinity,omitempty"`

	// Pod AntiAffinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod AntiAffinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAntiAffinity,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAntiAffinity *corev1.PodAntiAffinity `json:"podAntiAffinity,omitempty"`

	// Priority Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	PriorityClass string `json:"priorityClass,omitempty"`

	// Resource requirements
	// +optional
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:resourceRequirements,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.advanced,urn:alm:descriptor:com.tectonic.ui:advanced"
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Storage
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Storage SearchStorage `json:"storage,omitempty"`
}

type SearchStorage struct {
	// Data Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Data Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.storage"
	Data string `json:"data,omitempty"`

	// Extensions Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Extensions Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.storage"
	Extensions string `json:"extensions,omitempty"`

	// Storage Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Storage Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:StorageClass,urn:alm:descriptor:com.tectonic.ui:fieldGroup:node.search.storage"
	Class string `json:"class,omitempty"`
}

// SonarQubeStatus defines the observed state of SonarQube
type SonarQubeStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions status.Conditions `json:"conditions,omitempty"`

	// Kubernetes service that can be used to expose SonarQube
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Service string `json:"service,omitempty"`

	// Status of pods
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	Pods PodStatuses `json:"pods"`

	// Status of search pods
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Search Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	SearchPods PodStatuses `json:"searchPods"`

	// Status of instance
	Phase status.ConditionType `json:"phase,omitempty"`

	// Reason for status
	Reason string `json:"reason,omitempty"`

	// Expected revision for resources
	// Incremented when there is a change to spec or controller version
	Revision int32 `json:"revision,omitempty"`

	// Hash of latest spec & controller version for revision tracking
	RevisionHash string `json:"revisionHash,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQube is the Schema for the sonarqubes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sonarqubes,scope=Namespaced
// +operator-sdk:gen-csv:customresourcedefinitions.resources="StatefulSet,v1,\"\""
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
