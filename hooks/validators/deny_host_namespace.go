package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyHostNamespace is a Validator that denies sharing the host namespaces
type DenyHostNamespace struct{}

func (v DenyHostNamespace) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList
	if pod.Spec.HostNetwork {
		errs = append(errs, field.Forbidden(p.Child("hostNetwork"), "Host network is not allowed to be used"))
	}
	if pod.Spec.HostPID {
		errs = append(errs, field.Forbidden(p.Child("hostPID"), "Host pid is not allowed to be used"))
	}
	if pod.Spec.HostIPC {
		errs = append(errs, field.Forbidden(p.Child("hostIPC"), "Host ipc is not allowed to be used"))
	}
	return errs
}
