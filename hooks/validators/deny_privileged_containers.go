package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyPrivilegedContainers is a Validator that denies privileged containers
func DenyPrivilegedContainers(ctx context.Context, pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, c := range containers {
		if c.SecurityContext == nil || c.SecurityContext.Privileged == nil {
			continue
		}
		if *c.SecurityContext.Privileged == true {
			return admission.Denied("Privileged containers are not allowed")
		}
	}
	return admission.Allowed("ok")
}
