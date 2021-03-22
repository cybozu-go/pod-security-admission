package hooks

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podMutator struct {
	client  client.Client
	decoder *admission.Decoder
}

// NewPodMutator creates a webhook handler for Pod.
func NewPodMutator(c client.Client, dec *admission.Decoder, mutators []string) http.Handler {
	return &webhook.Admission{Handler: &podMutator{c, dec}}
}

func (m *podMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	po := &corev1.Pod{}
	err := m.decoder.Decode(req, po)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	poPatched := po.DeepCopy()

	marshaled, err := json.Marshal(poPatched)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}
