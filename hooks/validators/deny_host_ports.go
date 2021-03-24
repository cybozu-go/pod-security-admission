package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyHostPorts is a Validator that denies usage of HostPorts
func DenyHostPorts(ctx context.Context, pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.HostPort != 0 {
				return admission.Denied("Host port is not allowed to be used")
			}
		}
	}
	return admission.Allowed("ok")
}
