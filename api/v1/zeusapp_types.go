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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZeusappSpec defines the desired state of Zeusapp
type ZeusappSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Zeusapp. Edit zeusapp_types.go to remove/update
	Name          string   `json:"name,omitempty"`
	MinReplicas   int32    `json:"minReplicas"`
	Image         string   `json:"image"`
	Command       []string `json:"command"`
	ContainerPort int32    `json:"containerPort"`
}

// ZeusappStatus defines the observed state of Zeusapp
type ZeusappStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ReadyReplicas int32 `json:"readyReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Zeusapp is the Schema for the zeusapps API
type Zeusapp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZeusappSpec   `json:"spec,omitempty"`
	Status ZeusappStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ZeusappList contains a list of Zeusapp
type ZeusappList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zeusapp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zeusapp{}, &ZeusappList{})
}
