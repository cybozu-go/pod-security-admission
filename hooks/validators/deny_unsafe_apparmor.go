package validators

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyUnsafeAppArmor is a Validator that denies overriding or disabling the default AppArmor profile
type DenyUnsafeAppArmor struct{}

func (v DenyUnsafeAppArmor) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec").Child("annotations")
	var errs field.ErrorList

	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, corev1.DeprecatedAppArmorBetaContainerAnnotationKeyPrefix) &&
			v != corev1.DeprecatedAppArmorBetaProfileRuntimeDefault {
			errs = append(errs, field.Forbidden(p.Key(k), fmt.Sprintf("%s is not an allowed AppArmor profile", v)))
		}
	}

	p = field.NewPath("spec").Child("SecurityContext")
	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.AppArmorProfile != nil {
		errs = append(errs, validateAppArmorProfile(p, pod.Spec.SecurityContext.AppArmorProfile)...)
	}

	p = p.Child("containers")
	for i, co := range pod.Spec.Containers {
		if co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil {
			errs = append(errs, validateAppArmorProfile(p.Index(i), co.SecurityContext.AppArmorProfile)...)
		}
	}

	p = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		if co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil {
			errs = append(errs, validateAppArmorProfile(p.Index(i), co.SecurityContext.AppArmorProfile)...)
		}
	}

	p = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		if co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil {
			errs = append(errs, validateAppArmorProfile(p.Index(i), co.SecurityContext.AppArmorProfile)...)
		}
	}
	return errs
}

func validateAppArmorProfile(p *field.Path, profile *corev1.AppArmorProfile) field.ErrorList {
	if profile.Type == corev1.AppArmorProfileTypeRuntimeDefault {
		return nil
	}
	name := string(profile.Type)
	if profile.Type == corev1.AppArmorProfileTypeLocalhost && profile.LocalhostProfile != nil {
		name = fmt.Sprintf("%s/%s", profile.Type, *profile.LocalhostProfile)
	}
	return field.ErrorList{field.Forbidden(p, fmt.Sprintf("%s is not an allowed AppArmor profile", name))}
}
