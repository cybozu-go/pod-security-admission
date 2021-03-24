package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// default list of capabilities for Docker
// ref: https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities
var defaultCapabilities = []string{
	"AUDIT_WRITE",
	"CHOWN",
	"DAC_OVERRIDE",
	"FOWNER",
	"FSETID",
	"KILL",
	"MKNOD",
	"NET_BIND_SERVICE",
	"NET_RAW",
	"SETFCAP",
	"SETGID",
	"SETPCAP",
	"SETUID",
	"SYS_CHROOT",
}

// DenyUnsafeCapabilities is a Validator that denies adding capabilities beyond the default set
func DenyUnsafeCapabilities(ctx context.Context, pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, c := range containers {
		if c.SecurityContext == nil || c.SecurityContext.Capabilities == nil {
			continue
		}
		for _, add := range c.SecurityContext.Capabilities.Add {
			if !containsString(defaultCapabilities, string(add)) {
				return admission.Denied(fmt.Sprintf("Adding capability %s is not allowed", add))
			}
		}
	}

	return admission.Allowed("ok")
}
