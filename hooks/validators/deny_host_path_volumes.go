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
	ReadOnly   bool   `json:"readOnly"`
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

	for i, vol := range pod.Spec.Volumes {
		if vol.HostPath == nil {
			continue
		}
		allowed, readonly := v.allowedPath(vol.HostPath.Path)
		if !allowed {
			errs = append(errs, field.Forbidden(p.Index(i), "HostPath is not allowed to be used"))
			continue
		}

		if readonly {
			for _, container := range pod.Spec.Containers {
				for _, mount := range container.VolumeMounts {
					if mount.Name == vol.Name && !mount.ReadOnly {
						errs = append(errs, field.Forbidden(p.Index(i), "HostPath is allowed to be used only as read-only"))
					}
				}
			}
			for _, container := range pod.Spec.InitContainers {
				for _, mount := range container.VolumeMounts {
					if mount.Name == vol.Name && !mount.ReadOnly {
						errs = append(errs, field.Forbidden(p.Index(i), "HostPath is allowed to be used only as read-only"))
					}
				}
			}
			for _, container := range pod.Spec.EphemeralContainers {
				for _, mount := range container.VolumeMounts {
					if mount.Name == vol.Name && !mount.ReadOnly {
						errs = append(errs, field.Forbidden(p.Index(i), "HostPath is allowed to be used only as read-only"))
					}
				}
			}
		}
	}
	return errs
}

func (v DenyHostPathVolumes) allowedPath(path string) (bool, bool) {
	for _, allowed := range v.allowedHostPaths {
		if strings.HasPrefix(path, allowed.PathPrefix) && (len(path) == len(allowed.PathPrefix) || []rune(path)[len(allowed.PathPrefix)] == filepath.Separator) {
			return true, allowed.ReadOnly
		}
	}
	return false, false
}
