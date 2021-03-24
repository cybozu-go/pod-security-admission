package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var allowedSysctls = []string{
	"kernel.shm_rmid_forced",
	"net.ipv4.ip_local_port_range",
	"net.ipv4.tcp_syncookies",
	"net.ipv4.ping_group_range",
}

// DenyUnsafeSysctls is a Validator that denies usage of unsafe sysctls
func DenyUnsafeSysctls(ctx context.Context, pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext == nil {
		return admission.Allowed("ok")
	}
	for _, sysctl := range pod.Spec.SecurityContext.Sysctls {
		if !containsString(allowedSysctls, sysctl.Name) {
			return admission.Denied(fmt.Sprintf("Setting sysctl %s is not allowed", sysctl.Name))
		}
	}
	return admission.Allowed("ok")
}
