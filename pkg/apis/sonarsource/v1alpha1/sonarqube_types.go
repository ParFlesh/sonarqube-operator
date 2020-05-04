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
	// The name of a higher level application this instance is part of
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Application"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Application string `json:"application,omitempty"`

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

	// Secret with sonar configuration each key will be added as environment variables into Sonar Container.
	// (More Information: https://docs.sonarqube.org/latest/setup/environment-variables/)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Secret string `json:"secret,omitempty"`

	// Run SonarQube as a cluster
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Clustered"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Clustered bool `json:"clustered,omitempty"`

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

	// Application Node replicas (Only greater than 1 if clustered enabled)
	// (If set to 0, application and search nodes will be shutdown)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Replicas"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:number,urn:alm:descriptor:com.tectonic.ui:fieldGroup:instance"
	Replicas int32 `json:"replicas,omitempty"`

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

	// Status of pods
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Pods PodStatuses `json:"pods,omitempty"`

	// Status of search pods
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	SearchPods PodStatuses `json:"searchPods,omitempty"`

	// Pod Count
	Size int32 `json:"size,omitempty"`

	// Status of instance
	Phase status.ConditionType `json:"phase,omitempty"`

	// Reason for status
	Reason string `json:"reason,omitempty"`

	// Expected revision for resources
	// Incremented when there is a change to spec or controller version
	Revision int32 `json:"revision,omitempty"`

	// Hash of latest spec & controller version for revision tracking
	RevisionHash string `json:"revisionHas,omitempty"`
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
