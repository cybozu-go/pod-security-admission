package mutators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

type ForceRunAsNonRoot struct{}

func (m ForceRunAsNonRoot) Mutate(ctx context.Context, pod *corev1.Pod) bool {
	updated := false
	for i, co := range pod.Spec.Containers {
		sc := co.SecurityContext
		if sc == nil {
			sc = &corev1.SecurityContext{}
		}
		if sc.RunAsNonRoot == nil && sc.RunAsUser == nil {
			sc.RunAsNonRoot = pointer.BoolPtr(true)
			updated = true
		}
		pod.Spec.Containers[i].SecurityContext = sc
	}
	for i, co := range pod.Spec.InitContainers {
		sc := co.SecurityContext
		if sc == nil {
			sc = &corev1.SecurityContext{}
		}
		if sc.RunAsNonRoot == nil && sc.RunAsUser == nil {
			sc.RunAsNonRoot = pointer.BoolPtr(true)
			updated = true
		}
		pod.Spec.InitContainers[i].SecurityContext = sc
	}
	return updated
}
