package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateRunUser(p *field.Path, runAsNonRoot *bool, runAsUser *int64) *field.Error {
	if runAsNonRoot == nil && runAsUser == nil {
		return field.Forbidden(p, "RunAsNonRoot must be true")
	}
	if runAsNonRoot != nil && !*runAsNonRoot {
		return field.Forbidden(p.Child("runAsNonRoot"), "RunAsNonRoot must be true")
	}
	if runAsUser != nil && *runAsUser == 0 {
		return field.Forbidden(p.Child("runAsUser"), "Running with the root UID is forbidden")
	}
	return nil
}

// DenyRunAsRoot is a Validator that denies running as root users
type DenyRunAsRoot struct{}

func (v DenyRunAsRoot) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	allContainersAllowed := true
	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext == nil || (co.SecurityContext.RunAsNonRoot == nil && co.SecurityContext.RunAsUser == nil) {
			allContainersAllowed = false
			continue
		}
		if err := validateRunUser(pp.Index(i).Child("securityContext"), co.SecurityContext.RunAsNonRoot, co.SecurityContext.RunAsUser); err != nil {
			allContainersAllowed = false
			errs = append(errs, err)
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext == nil || (co.SecurityContext.RunAsNonRoot == nil && co.SecurityContext.RunAsUser == nil) {
			allContainersAllowed = false
			continue
		}
		if err := validateRunUser(pp.Index(i).Child("securityContext"), co.SecurityContext.RunAsNonRoot, co.SecurityContext.RunAsUser); err != nil {
			allContainersAllowed = false
			errs = append(errs, err)
		}
	}

	// runAsNonRoot of an ephemeral container will not be validated until the following issue is completed.
	// https://github.com/kubernetes/kubectl/issues/1108
	/*
		pp = p.Child("ephemeralContainers")
		for i, co := range pod.Spec.EphemeralContainers {
			if co.SecurityContext == nil || (co.SecurityContext.RunAsNonRoot == nil && co.SecurityContext.RunAsUser == nil) {
				allContainersAllowed = false
				continue
			}
			if err := validateRunUser(pp.Index(i).Child("securityContext"), co.SecurityContext.RunAsNonRoot, co.SecurityContext.RunAsUser); err != nil {
				allContainersAllowed = false
				errs = append(errs, err)
			}
		}
	*/

	if allContainersAllowed {
		return errs
	}

	if pod.Spec.SecurityContext == nil {
		errs = append(errs, field.Forbidden(p.Child("securityContext"), "RunAsNonRoot must be true"))
		return errs
	}
	if err := validateRunUser(p.Child("securityContext"), pod.Spec.SecurityContext.RunAsNonRoot, pod.Spec.SecurityContext.RunAsUser); err != nil {
		errs = append(errs, err)
	}
	return errs
}
