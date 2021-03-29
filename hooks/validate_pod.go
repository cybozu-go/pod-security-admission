package hooks

import (
	"context"
	"net/http"

	"github.com/cybozu-go/pod-security-admission/hooks/validators"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podValidator struct {
	client      client.Client
	log         logr.Logger
	decoder     *admission.Decoder
	profileName string
	validators  []validators.Validator
}

// NewPodValidator creates a webhook handler for Pod.
func NewPodValidator(c client.Client, log logr.Logger, dec *admission.Decoder, prof SecurityProfile) http.Handler {
	v := &podValidator{
		client:      c,
		log:         log,
		decoder:     dec,
		profileName: prof.Name,
		validators:  createValidators(prof),
	}
	return &webhook.Admission{Handler: v}
}

func createValidators(prof SecurityProfile) []validators.Validator {
	list := make([]validators.Validator, 0)
	if prof.DenyHostNamespace {
		list = append(list, validators.DenyHostNamespace{})
	}
	if prof.DenyPrivilegedContainers {
		list = append(list, validators.DenyPrivilegedContainers{})
	}
	if prof.DenyUnsafeCapabilities {
		list = append(list, validators.NewDenyUnsafeCapabilities(prof.AllowedCapabilities))
	}
	if prof.DenyHostPathVolumes {
		list = append(list, validators.DenyHostPathVolumes{})
	}
	if prof.DenyHostPorts {
		list = append(list, validators.NewDenyHostPorts(prof.AllowedHostPorts))
	}
	if prof.DenyUnsafeAppArmor {
		list = append(list, validators.DenyUnsafeAppArmor{})
	}
	if prof.DenyUnsafeSELinux {
		list = append(list, validators.DenyUnsafeSELinux{})
	}
	if prof.DenyUnsafeProcMount {
		list = append(list, validators.DenyUnsafeProcMount{})
	}
	if prof.DenyUnsafeSysctls {
		list = append(list, validators.DenyUnsafeSysctls{})
	}
	if prof.DenyNonCoreVolumeTypes {
		list = append(list, validators.DenyNonCoreVolumeTypes{})
	}
	if prof.DenyPrivilegeEscalation {
		list = append(list, validators.DenyPrivilegeEscalation{})
	}
	if prof.DenyRunAsRoot {
		list = append(list, validators.DenyRunAsRoot{})
	}
	if prof.DenyRootGroups {
		list = append(list, validators.DenyRootGroups{})
	}
	if prof.DenyUnsafeSeccomp {
		list = append(list, validators.DenyUnsafeSeccomp{})
	}
	return list
}

func (v *podValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	namespacedName := types.NamespacedName{
		Name:      req.Name,
		Namespace: req.Namespace,
	}
	v.log.Info("validating pod", "name", namespacedName, "profile", v.profileName)

	po := &corev1.Pod{}
	err := v.decoder.Decode(req, po)
	if err != nil {
		v.log.Error(err, "failed to decode pod", "name", namespacedName, "profile", v.profileName)
		return admission.Errored(http.StatusBadRequest, err)
	}

	var allErrs field.ErrorList
	for _, validator := range v.validators {
		errs := validator.Validate(ctx, po)
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) > 0 {
		reason := allErrs.ToAggregate().Error()
		v.log.Info("denied the pod", "name", namespacedName, "profile", v.profileName, "reason", reason)
		return admission.Denied(reason)
	}

	return admission.Allowed("ok")
}
