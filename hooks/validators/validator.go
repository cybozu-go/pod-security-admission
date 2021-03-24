package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Validator is a function signature for validators
type Validator func(ctx context.Context, pod *corev1.Pod) admission.Response

func containsString(list []string, item string) bool {
	for _, c := range list {
		if item == c {
			return true
		}
	}
	return false
}
