package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyHostPathVolumes is a Validator that denies usage of HostPath volumes
func DenyHostPathVolumes(ctx context.Context, pod *corev1.Pod) admission.Response {
	for _, vol := range pod.Spec.Volumes {
		if vol.HostPath == nil {
			continue
		}
		if len(vol.HostPath.Path) != 0 || vol.HostPath.Type != nil {
			return admission.Denied("Host path is not allowed to be used")
		}
	}
	return admission.Allowed("ok")
}
