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
	return errs
}
