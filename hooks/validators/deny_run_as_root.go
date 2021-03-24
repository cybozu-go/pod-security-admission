package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyRunAsRoot is a Validator that denies running as root users
func DenyRunAsRoot(ctx context.Context, pod *corev1.Pod) admission.Response {
	validate := func(runAsNonRoot *bool, runAsUser *int64) admission.Response {
		if runAsNonRoot == nil && runAsUser == nil {
			return admission.Denied("RunAsNonRoot must be true")
		}
		if runAsNonRoot != nil && *runAsNonRoot == false {
			return admission.Denied("RunAsNonRoot must be true")
		}
		if runAsUser != nil && *runAsUser == 0 {
			return admission.Denied("Running with the root UID is forbidden")
		}
		return admission.Allowed("ok")
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	allContainersAllowed := true
	for _, container := range containers {
		if container.SecurityContext == nil || (container.SecurityContext.RunAsNonRoot == nil && container.SecurityContext.RunAsUser == nil) {
			allContainersAllowed = false
			continue
		}
		if res := validate(container.SecurityContext.RunAsNonRoot, container.SecurityContext.RunAsUser); !res.Allowed {
			return res
		}
	}
	if allContainersAllowed {
		return admission.Allowed("ok")
	}

	if pod.Spec.SecurityContext == nil {
		return admission.Denied("RunAsNonRoot must be true")
	}
	return validate(pod.Spec.SecurityContext.RunAsNonRoot, pod.Spec.SecurityContext.RunAsUser)
}
