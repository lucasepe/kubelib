package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Reasons a resource is or is not ready.
const (
	ReasonAvailable   string = "Available"
	ReasonUnavailable string = "Unavailable"
	ReasonCreating    string = "Creating"
	ReasonDeleting    string = "Deleting"
)

// +kubebuilder:object:generate=true

// A Conditioned reflects the observed status of a resource.
// Only one condition of each type may exist.
type Conditioned struct {
	// Conditions of the resource.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// GetCondition returns the condition for the given ConditionType if exists,
// otherwise returns nil
func (s *Conditioned) GetCondition(ct string) metav1.Condition {
	for _, c := range s.Conditions {
		if c.Type == ct {
			return c
		}
	}

	return metav1.Condition{Type: ct, Status: metav1.ConditionUnknown}
}

// SetConditions sets the supplied conditions, replacing any existing conditions
// of the same type. This is a no-op if all supplied conditions are identical,
// ignoring the last transition time, to those already set.
func (s *Conditioned) SetConditions(c ...metav1.Condition) {
	for _, new := range c {
		exists := false
		for i, existing := range s.Conditions {
			if existing.Type != new.Type {
				continue
			}

			if conditionEqual(existing, new) {
				exists = true
				continue
			}

			s.Conditions[i] = new
			exists = true
		}
		if !exists {
			s.Conditions = append(s.Conditions, new)
		}
	}
}

// Creating returns a condition that indicates the resource is currently
// being created.
func Creating() metav1.Condition {
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonCreating,
	}
}

// Deleting returns a condition that indicates the resource is currently
// being deleted.
func Deleting() metav1.Condition {
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonDeleting,
	}
}

// Available returns a condition that indicates the resource is
// currently observed to be available for use.
func Available() metav1.Condition {
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonAvailable,
	}
}

// Unavailable returns a condition that indicates the resource is not
// currently available for use. Unavailable should be set only when Crossplane
// expects the resource to be available but knows it is not, for example
// because its API reports it is unhealthy.
func Unavailable() metav1.Condition {
	return metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ReasonUnavailable,
	}
}

// conditionEqual returns true if the condition is identical to the supplied condition,
// ignoring the LastTransitionTime.
func conditionEqual(c metav1.Condition, other metav1.Condition) bool {
	return c.Type == other.Type &&
		c.Status == other.Status &&
		c.Reason == other.Reason &&
		c.Message == other.Message
}
