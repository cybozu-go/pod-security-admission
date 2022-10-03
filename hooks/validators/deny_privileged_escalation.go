package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyPrivilegeEscalation is a Validator that denies privilege escalation
type DenyPrivilegeEscalation struct{}

func (v DenyPrivilegeEscalation) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || co.SecurityContext.AllowPrivilegeEscalation == nil {
			continue
		}
		escalation := *co.SecurityContext.AllowPrivilegeEscalation
		if escalation {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "allowPrivilegeEscalation"), "Allowing privilege escalation for containers is not allowed"))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || co.SecurityContext.AllowPrivilegeEscalation == nil {
			continue
		}
		escalation := *co.SecurityContext.AllowPrivilegeEscalation
		if escalation {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "allowPrivilegeEscalation"), "Allowing privilege escalation for containers is not allowed"))
		}
	}

	pp = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		if co.SecurityContext == nil || co.SecurityContext.AllowPrivilegeEscalation == nil {
			continue
		}
		escalation := *co.SecurityContext.AllowPrivilegeEscalation
		if escalation {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "allowPrivilegeEscalation"), "Allowing privilege escalation for containers is not allowed"))
		}
	}
	return errs
}
