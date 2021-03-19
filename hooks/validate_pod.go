package hooks

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-pod,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create,versions=v1,name=vpod.kb.io,admissionReviewVersions={v1,v1beta1}

type podValidator struct {
	client          client.Client
	decoder         *admission.Decoder
}

// NewPodValidator creates a webhook handler for Pod.
func NewPodValidator(c client.Client, dec *admission.Decoder) http.Handler {
	v := &podValidator{
		client:          c,
		decoder:         dec,
	}
	return &webhook.Admission{Handler: v}
}

func (v *podValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	po := &corev1.Pod{}
	err := v.decoder.Decode(req, po)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	containers := make([]corev1.Container, len(po.Spec.Containers)+len(po.Spec.InitContainers))
	copy(containers, po.Spec.Containers)
	copy(containers[len(po.Spec.Containers):], po.Spec.InitContainers)

	return admission.Allowed("ok")
}
