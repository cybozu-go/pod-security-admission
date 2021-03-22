package hooks

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create,versions=v1,name=mpod.kb.io,admissionReviewVersions={v1,v1beta1}

type podMutator struct {
	client  client.Client
	decoder *admission.Decoder
}

// NewPodMutator creates a webhook handler for Pod.
func NewPodMutator(c client.Client, dec *admission.Decoder) http.Handler {
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

