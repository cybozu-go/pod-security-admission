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

	p0 := field.NewPath("spec").Child("SecurityContext")
	hasPodAppArmorProfile := pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.AppArmorProfile != nil
	if hasPodAppArmorProfile {
		isTypeUnconfined := pod.Spec.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeUnconfined
		isTypeRuntimeDefault := pod.Spec.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
		isTypeLocalhost := pod.Spec.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
		hasNotAllowedType := !(isTypeUnconfined || isTypeRuntimeDefault || isTypeLocalhost)
		if hasNotAllowedType {
			errs = append(errs, field.Forbidden(p0.Child("AppArmorProfile"), "not an allowed *** AppArmor *** profile"))
		}
	}

	p1 := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		hasPodAppArmorProfile := co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil
		if hasPodAppArmorProfile {
			isTypeUnconfined := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeUnconfined
			isTypeRuntimeDefault := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
			isTypeLocalhost := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
			hasNotAllowedType := !(isTypeUnconfined || isTypeRuntimeDefault || isTypeLocalhost)
			if hasNotAllowedType {
				errs = append(errs, field.Forbidden(p1.Index(i), fmt.Sprintf("%s not an allowed *** AppArmor *** profile", "any")))
			}
		}
	}

	p2 := p.Child("initContainers")
	for i, co := range pod.Spec.Containers {
		hasPodAppArmorProfile := co.SecurityContext != nil && co.SecurityContext.AppArmorProfile != nil
		if hasPodAppArmorProfile {
			isTypeUnconfined := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeUnconfined
			isTypeRuntimeDefault := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeRuntimeDefault
			isTypeLocalhost := co.SecurityContext.AppArmorProfile.Type == corev1.AppArmorProfileTypeLocalhost
			hasNotAllowedType := !(isTypeUnconfined || isTypeRuntimeDefault || isTypeLocalhost)
			if hasNotAllowedType {
				errs = append(errs, field.Forbidden(p2.Index(i), fmt.Sprintf("%s not an allowed *** AppArmor *** profile", "any")))
			}
		}
	}
	return errs
}
