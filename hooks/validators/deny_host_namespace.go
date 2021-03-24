package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyHostNamespace is a Validator that denies sharing the host namespaces
func DenyHostNamespace(ctx context.Context, pod *corev1.Pod) admission.Response {
	if pod.Spec.HostNetwork {
		return admission.Denied("Host network is not allowed to be used")
	}
	if pod.Spec.HostPID {
		return admission.Denied("Host pid is not allowed to be used")
	}
	if pod.Spec.HostIPC {
		return admission.Denied("Host ipc is not allowed to be used")
	}
	return admission.Allowed("ok")
}
