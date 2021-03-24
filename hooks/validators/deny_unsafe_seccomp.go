package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyUnsafeSeccomp is a Validator that denies usage of non-default Seccomp profile
func DenyUnsafeSeccomp(ctx context.Context, pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SeccompProfile != nil && pod.Spec.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		return admission.Denied(fmt.Sprintf("%s is not an allowed seccomp prifile", pod.Spec.SecurityContext.SeccompProfile.Type))
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.SeccompProfile == nil {
			continue
		}
		if container.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			return admission.Denied(fmt.Sprintf("%s is not an allowed seccomp prifile", container.SecurityContext.SeccompProfile.Type))
		}
	}
	return admission.Allowed("ok")
}
