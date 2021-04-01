package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
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
type DenyUnsafeCapabilities struct {
	allowedCapabilities map[string]struct{}
}

func NewDenyUnsafeCapabilities(capabilities []string) DenyUnsafeCapabilities {
	allowedCapabilities := make(map[string]struct{})
	for _, c := range defaultCapabilities {
		allowedCapabilities[c] = struct{}{}
	}
	for _, c := range capabilities {
		allowedCapabilities[c] = struct{}{}
	}
	return DenyUnsafeCapabilities{allowedCapabilities: allowedCapabilities}
}

func (v DenyUnsafeCapabilities) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || co.SecurityContext.Capabilities == nil {
			continue
		}
		for j, add := range co.SecurityContext.Capabilities.Add {
			if _, ok := v.allowedCapabilities[string(add)]; !ok {
				errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "capabilities", "add").Index(j), fmt.Sprintf("Adding capability %s is not allowed", add)))
			}
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || co.SecurityContext.Capabilities == nil {
			continue
		}
		for j, add := range co.SecurityContext.Capabilities.Add {
			if _, ok := v.allowedCapabilities[string(add)]; !ok {
				errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "capabilities", "add").Index(j), fmt.Sprintf("Adding capability %s is not allowed", add)))
			}
		}
	}

	return errs
}
