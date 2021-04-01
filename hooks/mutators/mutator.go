package mutators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

// Mutator is an interface for mutator
type Mutator interface {
	// If the pod is mutated, it will return true.
	Mutate(ctx context.Context, pod *corev1.Pod) bool
}
