package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyUnsafeSELinux is a Validator that denies setting custom SELinux options
type DenyUnsafeSELinux struct{}

func (v DenyUnsafeSELinux) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SELinuxOptions != nil {
		errs = append(errs, field.Forbidden(p.Child("securityContext", "selinuxOptions"), "Setting custom SELinux options is not allowed"))
	}

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext != nil && co.SecurityContext.SELinuxOptions != nil {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "selinuxOptions"), "Setting custom SELinux options is not allowed"))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext != nil && co.SecurityContext.SELinuxOptions != nil {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "selinuxOptions"), "Setting custom SELinux options is not allowed"))
		}
	}

	return errs
}
