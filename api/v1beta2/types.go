package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeMemoryUsageLimitArgs defines the parameters for NodeMemoryUsageLimit plugin.
type NodeMemoryUsageLimitArgs struct {
	metav1.TypeMeta `json:",inline"`

	// NodeMemoryUsageLimit is the percentage limit (0-100]
	NodeMemoryUsageLimit int `json:"nodeMemoryUsageLimit,omitempty"`
}
