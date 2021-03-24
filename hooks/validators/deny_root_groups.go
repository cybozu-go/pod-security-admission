package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyRootGroups is a Validator that denies running with a root primary or supplementary GID
func DenyRootGroups(ctx context.Context, pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil {
		if pod.Spec.SecurityContext.RunAsGroup != nil && *pod.Spec.SecurityContext.RunAsGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
		for _, group := range pod.Spec.SecurityContext.SupplementalGroups {
			if group == 0 {
				return admission.Denied("Running with the supplementary GID is forbidden")
			}
		}
		if pod.Spec.SecurityContext.FSGroup != nil && *pod.Spec.SecurityContext.FSGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil {
			continue
		}
		if container.SecurityContext.RunAsGroup != nil && *container.SecurityContext.RunAsGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
	}
	return admission.Allowed("ok")
}
