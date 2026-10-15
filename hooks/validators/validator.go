package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validator is an interface for validator
type Validator interface {
	Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList
}

// usesUserNamespace returns true if the pod runs in its own user namespace (spec.hostUsers is false).
func usesUserNamespace(pod *corev1.Pod) bool {
	return pod.Spec.HostUsers != nil && !*pod.Spec.HostUsers
}
