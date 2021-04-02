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
	if !prof.HostNamespace {
		list = append(list, validators.DenyHostNamespace{})
	}
	if !prof.PrivilegedContainers {
		list = append(list, validators.DenyPrivilegedContainers{})
	}
	if !prof.Capabilities.UnsafeCapabilities {
		list = append(list, validators.NewDenyUnsafeCapabilities(prof.Capabilities.AllowedCapabilities))
	}
	if !prof.Volumes.HostPathVolumes {
		list = append(list, validators.DenyHostPathVolumes{})
	}
	if !prof.HostPorts.HostPorts {
		list = append(list, validators.NewDenyHostPorts(prof.HostPorts.AllowedHostPorts))
	}
	if !prof.UnsafeAppArmor {
		list = append(list, validators.DenyUnsafeAppArmor{})
	}
	if !prof.UnsafeSELinux {
		list = append(list, validators.DenyUnsafeSELinux{})
	}
	if !prof.UnsafeProcMount {
		list = append(list, validators.DenyUnsafeProcMount{})
	}
	if !prof.UnsafeSysctls {
		list = append(list, validators.DenyUnsafeSysctls{})
	}
	if !prof.Volumes.NonCoreVolumeTypes {
		list = append(list, validators.DenyNonCoreVolumeTypes{})
	}
	if !prof.PrivilegeEscalation {
		list = append(list, validators.DenyPrivilegeEscalation{})
	}
	if !prof.Users.RunAsRoot {
		list = append(list, validators.DenyRunAsRoot{})
	}
	if !prof.RootGroups {
		list = append(list, validators.DenyRootGroups{})
	}
	if !prof.UnsafeSeccomp {
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
