/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +kubebuilder:validation:Optional
package v1beta1

import (
	montymeta "github.com/aity-cloud/monty/pkg/util/meta"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type MontyClusterState string

const (
	MontyClusterStateError   MontyClusterState = "Error"
	MontyClusterStateWorking MontyClusterState = "Working"
	MontyClusterStateReady   MontyClusterState = "Ready"
)

// MontyClusterSpec defines the desired state of MontyCluster
type MontyClusterSpec struct {
	// +kubebuilder:default:=latest
	Version string `json:"version"`
	// +optional
	DefaultRepo          *string                         `json:"defaultRepo,omitempty"`
	NatsRef              corev1.LocalObjectReference     `json:"natsCluster"`
	Services             ServicesSpec                    `json:"services,omitempty"`
	Opensearch           *montymeta.OpensearchClusterRef `json:"opensearch,omitempty"`
	S3                   S3Spec                          `json:"s3,omitempty"`
	NulogHyperparameters map[string]intstr.IntOrString   `json:"nulogHyperparameters,omitempty"`
	DeployLogCollector   *bool                           `json:"deployLogCollector"`
	GlobalNodeSelector   map[string]string               `json:"globalNodeSelector,omitempty"`
	GlobalTolerations    []corev1.Toleration             `json:"globalTolerations,omitempty"`
}

// MontyClusterStatus defines the observed state of MontyCluster
type MontyClusterStatus struct {
	Conditions              []string          `json:"conditions,omitempty"`
	State                   MontyClusterState `json:"state,omitempty"`
	LogCollectorState       MontyClusterState `json:"logState,omitempty"`
	Auth                    AuthStatus        `json:"auth,omitempty"`
	PrometheusRuleNamespace string            `json:"prometheusRuleNamespace,omitempty"`
	IndexState              MontyClusterState `json:"indexState,omitempty"`
}

type AuthStatus struct {
	OpensearchAuthSecretKeyRef *corev1.SecretKeySelector `json:"opensearchAuthSecretKeyRef,omitempty"`
	S3Endpoint                 string                    `json:"s3Endpoint,omitempty"`
	S3AccessKey                *corev1.SecretKeySelector `json:"s3AccessKey,omitempty"`
	S3SecretKey                *corev1.SecretKeySelector `json:"s3SecretKey,omitempty"`
}

type OpensearchStatus struct {
	IndexState  MontyClusterState `json:"indexState,omitempty"`
	Version     *string           `json:"version,omitempty"`
	Initialized bool              `json:"initialized,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="State",type=boolean,JSONPath=`.status.state`
// +kubebuilder:storageversion

// MontyCluster is the Schema for the montyclusters API
type MontyCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MontyClusterSpec   `json:"spec,omitempty"`
	Status MontyClusterStatus `json:"status,omitempty"`
}

type ServicesSpec struct {
	Drain              DrainServiceSpec              `json:"drain,omitempty"`
	Inference          InferenceServiceSpec          `json:"inference,omitempty"`
	Preprocessing      PreprocessingServiceSpec      `json:"preprocessing,omitempty"`
	PayloadReceiver    PayloadReceiverServiceSpec    `json:"payloadReceiver,omitempty"`
	GPUController      GPUControllerServiceSpec      `json:"gpuController,omitempty"`
	Metrics            MetricsServiceSpec            `json:"metrics,omitempty"`
	OpensearchUpdate   OpensearchUpdateServiceSpec   `json:"opensearchUpdate,omitempty"`
	TrainingController TrainingControllerServiceSpec `json:"trainingController,omitempty"`
}

type DrainServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool                    `json:"enabled,omitempty"`
	NodeSelector        map[string]string        `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration      `json:"tolerations,omitempty"`
	Replicas            *int32                   `json:"replicas,omitempty"`
	Workload            WorkloadDrainServiceSpec `json:"workload,omitempty"`
}

type WorkloadDrainServiceSpec struct {
	Enabled  *bool  `json:"enabled,omitempty"`
	Replicas *int32 `json:"replicas,omitempty"`
}

type InferenceServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool                         `json:"enabled,omitempty"`
	PretrainedModels    []corev1.LocalObjectReference `json:"pretrainedModels,omitempty"`
	NodeSelector        map[string]string             `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration           `json:"tolerations,omitempty"`
}

type PreprocessingServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
	Replicas            *int32              `json:"replicas,omitempty"`
}

type PayloadReceiverServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type GPUControllerServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	RuntimeClass        *string             `json:"runtimeClass,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type TrainingControllerServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type MetricsServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool                          `json:"enabled,omitempty"`
	NodeSelector        map[string]string              `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration            `json:"tolerations,omitempty"`
	ExtraVolumeMounts   []montymeta.ExtraVolumeMount   `json:"extraVolumeMounts,omitempty"`
	PrometheusEndpoint  string                         `json:"prometheusEndpoint,omitempty"`
	PrometheusReference *montymeta.PrometheusReference `json:"prometheus,omitempty"`
}

type InsightsServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type UIServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type OpensearchUpdateServiceSpec struct {
	montymeta.ImageSpec `json:",inline,omitempty"`
	Enabled             *bool               `json:"enabled,omitempty"`
	NodeSelector        map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations         []corev1.Toleration `json:"tolerations,omitempty"`
}

type S3Spec struct {
	// If set, Monty will deploy an S3 pod to use internally.
	// Cannot be set at the same time as `external`.
	Internal *InternalSpec `json:"internal,omitempty"`
	// If set, Monty will connect to an external S3 endpoint.
	// Cannot be set at the same time as `internal`.
	External *ExternalSpec `json:"external,omitempty"`
	// Bucket used to persist nulog models.  If not set will use
	// monty-nulog-models.
	NulogS3Bucket string `json:"nulogS3Bucket,omitempty"`
	// Bucket used to persiste drain models.  It not set will use
	// monty-drain-models
	DrainS3Bucket string `json:"drainS3Bucket,omitempty"`
}

type InternalSpec struct {
	// Persistence configuration for internal S3 deployment. If unset, internal
	// S3 storage is not persistent.
	Persistence *montymeta.PersistenceSpec `json:"persistence,omitempty"`
}

type ExternalSpec struct {
	// +kubebuilder:validation:Required
	// External S3 endpoint URL.
	Endpoint string `json:"endpoint,omitempty"`
	// +kubebuilder:validation:Required
	// Reference to a secret containing "accessKey" and "secretKey" items. This
	// secret must already exist if specified.
	Credentials *corev1.SecretReference `json:"credentials,omitempty"`
}

// +kubebuilder:object:root=true

// MontyClusterList contains a list of MontyCluster
type MontyClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MontyCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MontyCluster{}, &MontyClusterList{})
}
