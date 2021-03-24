package hooks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cybozu-go/pod-security-admission/hooks/mutators"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podMutator struct {
	client       client.Client
	decoder      *admission.Decoder
	mutatorNames []string
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
	for _, name := range m.mutatorNames {
		mutator, ok := availableMutators[name]
		if !ok {
			return admission.Errored(http.StatusInternalServerError, errors.New("unknown mutator: "+name))
		}
		_ = mutator(ctx, poPatched) //NOTE: if there are multiple mutator, implement reinvocation with the return value.
	}
	marshaled, err := json.Marshal(poPatched)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}

var availableMutators = map[string]mutators.Mutator{
	"force-run-as-non-root": mutators.ForceRunAsNonRoot,
}
