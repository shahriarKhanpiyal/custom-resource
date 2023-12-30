/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CustomResourceSpec defines the desired state of CustomResource
type CustomResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	DeploymentName string `json:"deploymentName,omitempty"`
	//Replicas defines number of pods will be running in the deployment
	Replicas  *int32        `json:"replicas"`
	Container ContainerSpec `json:"container"`
	// Service contains ServiceName, ServiceType, ServiceNodePort
	//+optional
	Service ServiceSpec `json:"service,omitempty"`
}

type ContainerSpec struct {
	//+optional
	Image string `json:"image,omitempty"`
	//+optional
	Port int32 `json:"port,omitempty"`
}

type ServiceSpec struct {
	//+optional
	ServiceName string `json:"serviceName,omitempty"`
	ServiceType string `json:"serviceType"`
	//+optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
	ServicePort     int32 `json:"servicePort"`
}

// CustomResourceStatus defines the observed state of CustomResource
type CustomResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CustomResource is the Schema for the customresources API
type CustomResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomResourceSpec   `json:"spec,omitempty"`
	Status CustomResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CustomResourceList contains a list of CustomResource
type CustomResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomResource{}, &CustomResourceList{})
}
