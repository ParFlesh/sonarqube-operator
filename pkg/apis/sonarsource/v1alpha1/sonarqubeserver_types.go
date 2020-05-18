package v1alpha1

import (
	"github.com/operator-framework/operator-sdk/pkg/status"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SonarQubeServerSpec defines the desired state of SonarQubeServer
type SonarQubeServerSpec struct {
	// 0 and 1 are the only valid options.  Used to start and stop server.
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Size"
	// +kubebuilder:validation:Default=1
	Size int32 `json:"size,omitempty"`

	// Version of SonarQube image to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Version"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	// +kubebuilder:validation:Default=latest
	Version string `json:"version,omitempty"`

	// Image of SonarQube to deploy
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	// +kubebuilder:validation:Default=sonarqube
	Image string `json:"image,omitempty"`

	// Secret with sonar configuration files (sonar.properties).
	// Don't add cluster properties to configuration files as this could cause unexpected results
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	Secret string `json:"secret,omitempty"`

	// Sonar Node Type application or search when clustering is enabled otherwise aio (all-in-one)
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:select:aio,urn:alm:descriptor:com.tectonic.ui:select:application,urn:alm:descriptor:com.tectonic.ui:select:search"
	// +kubebuilder:validation:Default=aio
	// +kubebuilder:validation:Enum=aio;application;search
	Type ServerType `json:"type,omitempty"`

	// SonarQube application hosts list
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Hosts"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:hosts,urn:alm:descriptor:com.tectonic.ui:advanced"
	Hosts []string `json:"hosts,omitempty"`

	// SonarQube search hosts list
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Search Hosts"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch,urn:alm:descriptor:com.tectonic.ui:arrayFieldGroup:searchHosts,urn:alm:descriptor:com.tectonic.ui:advanced"
	SearchHosts []string `json:"searchHosts,omitempty"`

	// Deployment
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Deployment Deployment `json:"deployment,omitempty"`

	// Storage
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=false
	Storage Storage `json:"storage,omitempty"`
}

type Cluster struct {
}

type Deployment struct {
	// Node selector
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Selector"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:selector:Node,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Node Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Node Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:nodeAffinity,urn:alm:descriptor:com.tectonic.ui:advanced"
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`

	// Pod Affinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod Affinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAffinity,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAffinity *corev1.PodAffinity `json:"podAffinity,omitempty"`

	// Pod AntiAffinity
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Pod AntiAffinity"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podAntiAffinity,urn:alm:descriptor:com.tectonic.ui:advanced"
	PodAntiAffinity *corev1.PodAntiAffinity `json:"podAntiAffinity,omitempty"`

	// Priority Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:advanced"
	PriorityClass string `json:"priorityClass,omitempty"`

	// Resource requirements
	// +optional
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class"
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:resourceRequirements,urn:alm:descriptor:com.tectonic.ui:advanced"
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Service Account for running pods
	// +optional
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Service Account Name"
	//	+operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:selector:ServiceAccount,urn:alm:descriptor:com.tectonic.ui:advanced"
	// +kubebuilder:validation:Default=default
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

type Storage struct {
	// Data Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Data Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:storage"
	// +kubebuilder:validation:Default=1Gi
	DataSize string `json:"dataSize,omitempty"`

	// Storage Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Data Storage Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:StorageClass,urn:alm:descriptor:com.tectonic.ui:fieldGroup:storage"
	DataClass *string `json:"dataClass,omitempty"`

	// Extensions Volume Size
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Extensions Volume Size"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text,urn:alm:descriptor:com.tectonic.ui:fieldGroup:storage"
	// +kubebuilder:validation:Default=1Gi
	ExtensionsSize string `json:"extensionsSize,omitempty"`

	// Storage Class Name
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Extensions Storage Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:StorageClass,urn:alm:descriptor:com.tectonic.ui:fieldGroup:storage"
	ExtensionsClass *string `json:"extensionsClass,omitempty"`
}

// SonarQubeServerStatus defines the observed state of SonarQubeServer
type SonarQubeServerStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions status.Conditions `json:"conditions,omitempty"`

	// Kubernetes service that can be used to expose SonarQubeServer
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Service"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Service string `json:"service,omitempty"`

	// Status of pods
	// +optional
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Pod Statuses"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	// +kubebuilder:validation:Default={Available:[],Progressing:[],ReplicaFailure:[]}
	Deployment DeploymentStatus `json:"deployment"`

	// Expected revision for resources
	// Incremented when there is a change to spec or controller version
	// +kubebuilder:validation:Default=0
	Revision int32 `json:"revision,omitempty"`

	// Hash of latest spec & controller version for revision tracking
	RevisionHash string `json:"revisionHash,omitempty"`
}

type DeploymentStatus map[appsv1.DeploymentConditionType][]string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQubeServer is the Schema for the sonarqubeservers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sonarqubeservers,scope=Namespaced
type SonarQubeServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SonarQubeServerSpec   `json:"spec,omitempty"`
	Status SonarQubeServerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SonarQubeServerList contains a list of SonarQubeServer
type SonarQubeServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SonarQubeServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SonarQubeServer{}, &SonarQubeServerList{})
}
