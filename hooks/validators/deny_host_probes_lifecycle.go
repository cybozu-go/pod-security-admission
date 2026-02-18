package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyHostProbesAndLifecycle is a Validator that denies setting host field
// in probes and lifecycle handlers.
type DenyHostProbesAndLifecycle struct{}

func (v DenyHostProbesAndLifecycle) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		errs = append(errs, validateContainerProbesAndLifecycle(pp.Index(i), &co)...)
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		errs = append(errs, validateContainerProbesAndLifecycle(pp.Index(i), &co)...)
	}

	pp = p.Child("ephemeralContainers")
	for i, co := range pod.Spec.EphemeralContainers {
		errs = append(errs, validateEphemeralContainerProbesAndLifecycle(pp.Index(i), &co)...)
	}

	return errs
}

func validateContainerProbesAndLifecycle(p *field.Path, co *corev1.Container) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, validateProbeHost(p.Child("livenessProbe"), co.LivenessProbe)...)
	errs = append(errs, validateProbeHost(p.Child("readinessProbe"), co.ReadinessProbe)...)
	errs = append(errs, validateProbeHost(p.Child("startupProbe"), co.StartupProbe)...)
	if co.Lifecycle != nil {
		errs = append(errs, validateLifecycleHandlerHost(p.Child("lifecycle", "postStart"), co.Lifecycle.PostStart)...)
		errs = append(errs, validateLifecycleHandlerHost(p.Child("lifecycle", "preStop"), co.Lifecycle.PreStop)...)
	}
	return errs
}

func validateEphemeralContainerProbesAndLifecycle(p *field.Path, co *corev1.EphemeralContainer) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, validateProbeHost(p.Child("livenessProbe"), co.LivenessProbe)...)
	errs = append(errs, validateProbeHost(p.Child("readinessProbe"), co.ReadinessProbe)...)
	errs = append(errs, validateProbeHost(p.Child("startupProbe"), co.StartupProbe)...)
	if co.Lifecycle != nil {
		errs = append(errs, validateLifecycleHandlerHost(p.Child("lifecycle", "postStart"), co.Lifecycle.PostStart)...)
		errs = append(errs, validateLifecycleHandlerHost(p.Child("lifecycle", "preStop"), co.Lifecycle.PreStop)...)
	}
	return errs
}

func validateProbeHost(p *field.Path, probe *corev1.Probe) field.ErrorList {
	if probe == nil {
		return nil
	}
	var errs field.ErrorList
	if probe.HTTPGet != nil && probe.HTTPGet.Host != "" {
		errs = append(errs, field.Forbidden(p.Child("httpGet", "host"), "Host is not allowed to be set"))
	}
	if probe.TCPSocket != nil && probe.TCPSocket.Host != "" {
		errs = append(errs, field.Forbidden(p.Child("tcpSocket", "host"), "Host is not allowed to be set"))
	}
	return errs
}

func validateLifecycleHandlerHost(p *field.Path, handler *corev1.LifecycleHandler) field.ErrorList {
	if handler == nil {
		return nil
	}
	var errs field.ErrorList
	if handler.HTTPGet != nil && handler.HTTPGet.Host != "" {
		errs = append(errs, field.Forbidden(p.Child("httpGet", "host"), "Host is not allowed to be set"))
	}
	if handler.TCPSocket != nil && handler.TCPSocket.Host != "" {
		errs = append(errs, field.Forbidden(p.Child("tcpSocket", "host"), "Host is not allowed to be set"))
	}
	return errs
}
