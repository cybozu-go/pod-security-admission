package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyUnsafeSeccomp is a Validator that denies usage of non-default Seccomp profile
type DenyUnsafeSeccomp struct{}

func (v DenyUnsafeSeccomp) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SeccompProfile != nil && pod.Spec.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		errs = append(errs, field.Forbidden(p.Child("securityContext", "seccompProfile", "type"), fmt.Sprintf("%s is not an allowed seccomp profile", pod.Spec.SecurityContext.SeccompProfile.Type)))
	}

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || co.SecurityContext.SeccompProfile == nil {
			continue
		}
		if co.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "seccompProfile", "type"), fmt.Sprintf("%s is not an allowed seccomp profile", co.SecurityContext.SeccompProfile.Type)))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || co.SecurityContext.SeccompProfile == nil {
			continue
		}
		if co.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "seccompProfile", "type"), fmt.Sprintf("%s is not an allowed seccomp profile", co.SecurityContext.SeccompProfile.Type)))
		}
	}

	return errs
}
