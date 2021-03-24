package hooks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podMutator struct {
	client   client.Client
	decoder  *admission.Decoder
	mutators []string
}

// NewPodMutator creates a webhook handler for Pod.
func NewPodMutator(c client.Client, dec *admission.Decoder, mutators []string) http.Handler {
	return &webhook.Admission{Handler: &podMutator{c, dec, mutators}}
}

func (m *podMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	po := &corev1.Pod{}
	err := m.decoder.Decode(req, po)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	poPatched := po.DeepCopy()

REINVOCATION:
	for _, name := range m.mutators {
		mutator, ok := mutators[name]
		if !ok {
			return admission.Errored(http.StatusInternalServerError, errors.New("unknown mutator: "+name))
		}
		updated := mutator(poPatched)
		if updated {
			continue REINVOCATION
		}
	}
	marshaled, err := json.Marshal(poPatched)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}

type Mutator func(pod *corev1.Pod) bool

var mutators = map[string]Mutator{
	"force-run-as-non-root": forceRunAsNonRoot,
}

func forceRunAsNonRoot(pod *corev1.Pod) bool {
	updated := false
	for i, co := range pod.Spec.Containers {
		sc := co.SecurityContext
		if sc == nil {
			sc = &corev1.SecurityContext{}
		}
		if sc.RunAsNonRoot == nil && sc.RunAsUser == nil {
			sc.RunAsNonRoot = pointer.BoolPtr(true)
			updated = true
		}
		pod.Spec.Containers[i].SecurityContext = sc
	}
	for i, co := range pod.Spec.InitContainers {
		sc := co.SecurityContext
		if sc == nil {
			sc = &corev1.SecurityContext{}
		}
		if sc.RunAsNonRoot == nil && sc.RunAsUser == nil {
			sc.RunAsNonRoot = pointer.BoolPtr(true)
			updated = true
		}
		pod.Spec.InitContainers[i].SecurityContext = sc
	}
	return updated
}
