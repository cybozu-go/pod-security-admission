package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var allowedSysctls = map[string]struct{}{
	"kernel.shm_rmid_forced":              {},
	"net.ipv4.ip_local_port_range":        {},
	"net.ipv4.tcp_syncookies":             {},
	"net.ipv4.ping_group_range":           {},
	"net.ipv4.ip_unprivileged_port_start": {},
	"net.ipv4.ip_local_reserved_ports":    {}, // since Kubernetes 1.27
	"net.ipv4.tcp_keepalive_time":         {}, // since Kubernetes 1.29
	"net.ipv4.tcp_fin_timeout":            {}, // since Kubernetes 1.29
	"net.ipv4.tcp_keepalive_intvl":        {}, // since Kubernetes 1.29
	"net.ipv4.tcp_keepalive_probes":       {}, // since Kubernetes 1.29
}

// DenyUnsafeSysctls is a Validator that denies usage of unsafe sysctls
type DenyUnsafeSysctls struct{}

func (v DenyUnsafeSysctls) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	if pod.Spec.SecurityContext == nil {
		return errs
	}
	for i, sysctl := range pod.Spec.SecurityContext.Sysctls {
		if _, ok := allowedSysctls[sysctl.Name]; !ok {
			errs = append(errs, field.Forbidden(p.Child("securityContext", "sysctls").Index(i), fmt.Sprintf("Setting sysctl %s is not allowed", sysctl.Name)))
		}
	}
	return errs
}
