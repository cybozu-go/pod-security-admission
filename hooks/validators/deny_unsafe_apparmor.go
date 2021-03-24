package validators

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyUnsafeAppArmor is a Validator that denies overriding or disabling the default AppArmor profile
func DenyUnsafeAppArmor(ctx context.Context, pod *corev1.Pod) admission.Response {
	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, corev1.AppArmorBetaContainerAnnotationKeyPrefix) &&
			v != corev1.AppArmorBetaProfileRuntimeDefault {
			return admission.Denied(fmt.Sprintf("%s is not an allowed AppArmor profile", v))
		}
	}
	return admission.Allowed("ok")
}
