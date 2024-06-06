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
	hasPodAppArmorProfile := pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.AppArmorProfile != nil
	if hasPodAppArmorProfile {
		isTypeRuntimeDefault := pod.Spec.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
		isTypeLocalhost := pod.Spec.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
		hasNotAllowedType := !(isTypeRuntimeDefault || isTypeLocalhost)
		if hasNotAllowedType {
			errs = append(errs, field.Forbidden(p, fmt.Sprintf("%v is not an allowed AppArmor profile", pod.Spec.SecurityContext.AppArmorProfile.Type)))
		}
	}

	p = p.Child("containers")
	for i, co := range pod.Spec.Containers {
		hasPodAppArmorProfile := co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil
		if hasPodAppArmorProfile {
			isTypeRuntimeDefault := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
			isTypeLocalhost := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
			hasNotAllowedType := !(isTypeRuntimeDefault || isTypeLocalhost)
			if hasNotAllowedType {
				errs = append(errs, field.Forbidden(p.Index(i), fmt.Sprintf("%v is not an allowed AppArmor profile", co.SecurityContext.AppArmorProfile.Type)))
			}
		}
	}

	p = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		hasPodAppArmorProfile := co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil
		if hasPodAppArmorProfile {
			isTypeRuntimeDefault := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
			isTypeLocalhost := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
			hasNotAllowedType := !(isTypeRuntimeDefault || isTypeLocalhost)
			if hasNotAllowedType {
				errs = append(errs, field.Forbidden(p.Index(i), fmt.Sprintf("%v is not an allowed AppArmor profile", co.SecurityContext.AppArmorProfile.Type)))
			}
		}
	}

	p = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		hasPodAppArmorProfile := co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil
		if hasPodAppArmorProfile {
			isTypeRuntimeDefault := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
			isTypeLocalhost := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
			hasNotAllowedType := !(isTypeRuntimeDefault || isTypeLocalhost)
			if hasNotAllowedType {
				errs = append(errs, field.Forbidden(p.Index(i), fmt.Sprintf("%v is not an allowed AppArmor profile", co.SecurityContext.AppArmorProfile.Type)))
			}
		}
	}
	return errs
}
