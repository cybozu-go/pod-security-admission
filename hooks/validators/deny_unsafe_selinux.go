package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyUnsafeSELinux is a Validator that denies setting custom SELinux options
func DenyUnsafeSELinux(ctx context.Context, pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SELinuxOptions != nil {
		return admission.Denied("Setting custom SELinux options is not allowed")
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext != nil && container.SecurityContext.SELinuxOptions != nil {
			return admission.Denied("Setting custom SELinux options is not allowed")
		}
	}
	return admission.Allowed("ok")
}
