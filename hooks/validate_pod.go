package hooks

import (
	"context"
	"errors"
	"net/http"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/cybozu-go/pod-security-admission/hooks/validators"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podValidator struct {
	client         client.Client
	decoder        *admission.Decoder
	validatorNames []string
}

// NewPodValidator creates a webhook handler for Pod.
func NewPodValidator(c client.Client, dec *admission.Decoder, validators []string) http.Handler {
	v := &podValidator{
		client:         c,
		decoder:        dec,
		validatorNames: validators,
	}
	return &webhook.Admission{Handler: v}
}

func (v *podValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	po := &corev1.Pod{}
	err := v.decoder.Decode(req, po)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	var allErrs field.ErrorList
	for _, name := range v.validatorNames {
		validator, ok := availableValidators[name]
		if !ok {
			return admission.Errored(http.StatusInternalServerError, errors.New("unknown validator: "+name))
		}
		errs := validator(ctx, po)
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) > 0 {
		return admission.Denied(allErrs.ToAggregate().Error())
	}

	return admission.Allowed("ok")
}

var availableValidators = map[string]validators.Validator{
	"deny-host-namespace":        validators.DenyHostNamespace,
	"deny-privileged-containers": validators.DenyPrivilegedContainers,
	"deny-unsafe-capabilities":   validators.DenyUnsafeCapabilities,
	"deny-host-path-volumes":     validators.DenyHostPathVolumes,
	"deny-host-ports":            validators.DenyHostPorts,
	"deny-unsafe-apparmor":       validators.DenyUnsafeAppArmor,
	"deny-unsafe-selinux":        validators.DenyUnsafeSELinux,
	"deny-unsafe-proc-mount":     validators.DenyUnsafeProcMount,
	"deny-unsafe-sysctls":        validators.DenyUnsafeSysctls,
	"deny-non-core-volume-types": validators.DenyNonCoreVolumeTypes,
	"deny-privilege-escalation":  validators.DenyPrivilegeEscalation,
	"deny-run-as-root":           validators.DenyRunAsRoot,
	"deny-root-groups":           validators.DenyRootGroups,
	"deny-unsafe-seccomp":        validators.DenyUnsafeSeccomp,
}
