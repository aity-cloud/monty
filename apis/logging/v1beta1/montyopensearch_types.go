package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	opsterv1 "opensearch.opster.io/api/v1"
)

type MontyOpensearchState string

const (
	MontyOpensearchStateError   MontyOpensearchState = "Error"
	MontyOpensearchStateWorking MontyOpensearchState = "Working"
	MontyOpensearchStateReady   MontyOpensearchState = "Ready"
)

type OpensearchS3Protocol string

const (
	OpensearchS3ProtocolHTTPS OpensearchS3Protocol = "https"
	OpensearchS3ProtocolHTTP  OpensearchS3Protocol = "http"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="State",type=boolean,JSONPath=`.status.state`

type MontyOpensearch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MontyOpensearchSpec   `json:"spec,omitempty"`
	Status MontyOpensearchStatus `json:"status,omitempty"`
}

type MontyOpensearchStatus struct {
	Conditions        []string             `json:"conditions,omitempty"`
	State             MontyOpensearchState `json:"state,omitempty"`
	OpensearchVersion *string              `json:"opensearchVersion,omitempty"`
	Version           *string              `json:"version,omitempty"`
	PasswordGenerated bool                 `json:"passwordGenerated,omitempty"`
}

type MontyOpensearchSpec struct {
	*ClusterConfigSpec `json:",inline"`
	OpensearchSettings `json:"opensearch,omitempty"`
	ExternalURL        string                       `json:"externalURL,omitempty"`
	ImageRepo          string                       `json:"imageRepo"`
	OpensearchVersion  string                       `json:"opensearchVersion,omitempty"`
	Version            string                       `json:"version,omitempty"`
	NatsRef            *corev1.LocalObjectReference `json:"natsCluster"`
}

type OpensearchSettings struct {
	ImageOverride *string                   `json:"imageOverride,omitempty"`
	NodePools     []opsterv1.NodePool       `json:"nodePools,omitempty"`
	Dashboards    opsterv1.DashboardsConfig `json:"dashboards,omitempty"`
	Security      *opsterv1.Security        `json:"security,omitempty"`
	S3Settings    *OpensearchS3Settings     `json:"s3,omitempty"`
}

type OpensearchS3Settings struct {
	Endpoint         string                      `json:"endpoint,omitempty"`
	PathStyleAccess  bool                        `json:"pathStyleAccess,omitempty"`
	Protocol         OpensearchS3Protocol        `json:"protocol,omitempty"`
	ProxyHost        string                      `json:"proxyHost,omitempty"`
	ProxyPort        *int32                      `json:"proxyPort,omitempty"`
	CredentialSecret corev1.LocalObjectReference `json:"credentialSecret,omitempty"`
	Repository       S3PathSettings              `json:"repository,omitempty"`
}

// +kubebuilder:object:root=true

// MontyOpensearchList contains a list of MontyOpensearch
type MontyOpensearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MontyOpensearch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MontyOpensearch{}, &MontyOpensearchList{})
}
