package validators

import (
	"context"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type AllowedHostPath struct {
	PathPrefix string `json:"pathPrefix"`
}

// DenyHostPathVolumes is a Validator that denies usage of HostPath volumes
type DenyHostPathVolumes struct {
	allowedHostPaths []AllowedHostPath
}

func NewDenyHostPaths(hostPaths []AllowedHostPath) DenyHostPathVolumes {
	return DenyHostPathVolumes{allowedHostPaths: hostPaths}
}

func (v DenyHostPathVolumes) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec").Child("volumes")
	var errs field.ErrorList

volume_for:
	for i, vol := range pod.Spec.Volumes {
		if vol.HostPath == nil {
			continue
		}
		if v.allowedHostPaths == nil {
			errs = append(errs, field.Forbidden(p.Index(i), "HostPath is not allowed to be used"))
		} else {
			pathstr := vol.HostPath.Path
			for _, allowed := range v.allowedHostPaths {
				if strings.HasPrefix(pathstr, allowed.PathPrefix) && (len(pathstr) == len(allowed.PathPrefix) || []rune(pathstr)[len(allowed.PathPrefix)] == filepath.Separator) {
					continue volume_for
				}
			}
			errs = append(errs, field.Forbidden(p.Index(i), "HostPath is not allowed to be used"))
		}
	}
	return errs
}
