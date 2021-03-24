package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyUnsafeProcMount is a Validator that denies unmasked proc mount
func DenyUnsafeProcMount(ctx context.Context, pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.ProcMount == nil {
			continue
		}
		proc := *container.SecurityContext.ProcMount
		if proc != corev1.DefaultProcMount {
			return admission.Denied(fmt.Sprintf("ProcMountType %s is not allowed", proc))
		}
	}
	return admission.Allowed("ok")
}
