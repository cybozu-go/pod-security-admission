package validators

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyUnsafeProcMount is a Validator that denies unmasked proc mount
type DenyUnsafeProcMount struct{}

func (v DenyUnsafeProcMount) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || co.SecurityContext.ProcMount == nil {
			continue
		}
		proc := *co.SecurityContext.ProcMount
		if proc != corev1.DefaultProcMount {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "procMount"), fmt.Sprintf("ProcMountType %s is not allowed", proc)))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || co.SecurityContext.ProcMount == nil {
			continue
		}
		proc := *co.SecurityContext.ProcMount
		if proc != corev1.DefaultProcMount {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "procMount"), fmt.Sprintf("ProcMountType %s is not allowed", proc)))
		}
	}

	pp = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		if co.SecurityContext == nil || co.SecurityContext.ProcMount == nil {
			continue
		}
		proc := *co.SecurityContext.ProcMount
		if proc != corev1.DefaultProcMount {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "procMount"), fmt.Sprintf("ProcMountType %s is not allowed", proc)))
		}
	}

	return errs
}
