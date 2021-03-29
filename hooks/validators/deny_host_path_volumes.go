package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyHostPathVolumes is a Validator that denies usage of HostPath volumes
type DenyHostPathVolumes struct{}

func (v DenyHostPathVolumes) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec").Child("volumes")
	var errs field.ErrorList
	for i, vol := range pod.Spec.Volumes {
		if vol.HostPath == nil {
			continue
		}
		if len(vol.HostPath.Path) != 0 || vol.HostPath.Type != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type HostPath is not allowed to be used"))
		}
	}
	return errs
}
