package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyPrivilegeEscalation is a Validator that denies privilege escalation
func DenyPrivilegeEscalation(ctx context.Context, pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.AllowPrivilegeEscalation == nil {
			continue
		}
		escalation := *container.SecurityContext.AllowPrivilegeEscalation
		if escalation {
			return admission.Denied("Allowing privilege escalation for containers is not allowed")
		}
	}
	return admission.Allowed("ok")
}
