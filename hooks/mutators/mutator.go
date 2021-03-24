package mutators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

// Mutator is a function signature for mutators
// If the pod is mutated, return true.
type Mutator func(ctx context.Context, pod *corev1.Pod) bool
