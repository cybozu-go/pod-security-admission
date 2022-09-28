package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyRootGroups is a Validator that denies running with a root primary or supplementary GID
type DenyRootGroups struct{}

func (v DenyRootGroups) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	if pod.Spec.SecurityContext != nil {
		pp := p.Child("securityContext")
		if pod.Spec.SecurityContext.RunAsGroup != nil && *pod.Spec.SecurityContext.RunAsGroup == 0 {
			errs = append(errs, field.Forbidden(pp.Child("runAsGroup"), "Running with the root GID is forbidden"))
		}
		for i, group := range pod.Spec.SecurityContext.SupplementalGroups {
			if group == 0 {
				errs = append(errs, field.Forbidden(pp.Child("supplementalGroups").Index(i), "Running with the supplementary GID is forbidden"))
			}
		}
		if pod.Spec.SecurityContext.FSGroup != nil && *pod.Spec.SecurityContext.FSGroup == 0 {
			errs = append(errs, field.Forbidden(pp.Child("fsGroup"), "Running with the root GID is forbidden"))
		}
	}

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil {
			continue
		}
		if co.SecurityContext.RunAsGroup != nil && *co.SecurityContext.RunAsGroup == 0 {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "runAsGroup"), "Running with the root GID is forbidden"))
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil {
			continue
		}
		if co.SecurityContext.RunAsGroup != nil && *co.SecurityContext.RunAsGroup == 0 {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "runAsGroup"), "Running with the root GID is forbidden"))
		}
	}

	pp = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		if co.SecurityContext == nil {
			continue
		}
		if co.SecurityContext.RunAsGroup != nil && *co.SecurityContext.RunAsGroup == 0 {
			errs = append(errs, field.Forbidden(pp.Index(i).Child("securityContext", "runAsGroup"), "Running with the root GID is forbidden"))
		}
	}
	return errs
}
