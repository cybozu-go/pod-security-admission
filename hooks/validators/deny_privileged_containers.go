package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyPrivilegedContainers is a Validator that denies privileged containers
type DenyPrivilegedContainers struct{}

func (v DenyPrivilegedContainers) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || co.SecurityContext.Privileged == nil {
			continue
		}
		if *co.SecurityContext.Privileged {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "privileged"), "Privileged containers are not allowed"))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || co.SecurityContext.Privileged == nil {
			continue
		}
		if *co.SecurityContext.Privileged {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "privileged"), "Privileged containers are not allowed"))
		}
	}

	pp = p.Child("ephemeralContainer")
	for i, co := range pod.Spec.EphemeralContainers {
		if co.SecurityContext == nil || co.SecurityContext.Privileged == nil {
			continue
		}
		if *co.SecurityContext.Privileged {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "privileged"), "Privileged containers are not allowed"))
		}
	}
	return errs
}
