package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var allowedSysctls = []string{
	"kernel.shm_rmid_forced",
	"net.ipv4.ip_local_port_range",
	"net.ipv4.tcp_syncookies",
	"net.ipv4.ping_group_range",
}

// DenyUnsafeSysctls is a Validator that denies usage of unsafe sysctls
func DenyUnsafeSysctls(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	if pod.Spec.SecurityContext == nil {
		return errs
	}
	for i, sysctl := range pod.Spec.SecurityContext.Sysctls {
		if !containsString(allowedSysctls, sysctl.Name) {
			errs = append(errs, field.Forbidden(p.Index(i).Child("securityContext", "sysctls"), fmt.Sprintf("Setting sysctl %s is not allowed", sysctl.Name)))
		}
	}
	return errs
}
