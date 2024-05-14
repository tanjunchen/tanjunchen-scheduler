package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DynamicArgs is the args struction of scheduler plugin.
type DynamicArgs struct {
	metav1.TypeMeta     `json:",inline"`
	ToleranceCPURate    float64 `json:"toleranceCPURate,omitempty"`
	ToleranceMemoryRate float64 `json:"toleranceMemoryRate,omitempty"`
}
