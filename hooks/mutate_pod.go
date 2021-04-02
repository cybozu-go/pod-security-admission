package hooks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cybozu-go/pod-security-admission/hooks/mutators"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podMutator struct {
	client      client.Client
	log         logr.Logger
	decoder     *admission.Decoder
	profileName string
	mutators    []mutators.Mutator
}

// NewPodMutator creates a webhook handler for Pod.
func NewPodMutator(c client.Client, log logr.Logger, dec *admission.Decoder, prof SecurityProfile) http.Handler {
	m := &podMutator{
		client:      c,
		log:         log,
		decoder:     dec,
		profileName: prof.Name,
		mutators:    createMutators(prof),
	}
	return &webhook.Admission{Handler: m}
}

func createMutators(prof SecurityProfile) []mutators.Mutator {
	list := make([]mutators.Mutator, 0)
	if prof.ForceRunAsNonRoot {
		list = append(list, mutators.ForceRunAsNonRoot{})
	}
	return list
}

func (m *podMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	namespacedName := types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}
	m.log.Info("mutating pod,", "name", namespacedName, "profile", m.profileName)

	po := &corev1.Pod{}
	err := m.decoder.Decode(req, po)
	if err != nil {
		m.log.Error(err, "failed to decode pod", "name", namespacedName, "profile", m.profileName)
		return admission.Errored(http.StatusBadRequest, err)
	}

	poPatched := po.DeepCopy()
	for _, mutator := range m.mutators {
		mutated := mutator.Mutate(ctx, poPatched)
		if mutated {
			m.log.Info("pod mutated", "name", namespacedName, "profile", m.profileName)
		}
	}
	marshaled, err := json.Marshal(poPatched)
	if err != nil {
		m.log.Error(err, "failed to marshal patch", "name", namespacedName, "profile", m.profileName)
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}
