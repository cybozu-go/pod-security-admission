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
