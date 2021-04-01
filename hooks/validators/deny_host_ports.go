package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type PortRange struct {
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}

// DenyHostPorts is a Validator that denies usage of HostPorts
type DenyHostPorts struct {
	allowedHostPorts []PortRange
}

func NewDenyHostPorts(hostPorts []PortRange) DenyHostPorts {
	return DenyHostPorts{allowedHostPorts: hostPorts}
}

func (v DenyHostPorts) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		ppp := pp.Index(i)
		for j, port := range co.Ports {
			allowed := false
			if port.HostPort == 0 {
				allowed = true
			}
			for _, r := range v.allowedHostPorts {
				if r.Min <= port.HostPort && port.HostPort <= r.Max {
					allowed = true
				}
			}
			if !allowed {
				errs = append(errs, field.Forbidden(ppp.Child("ports").Index(j), "Host port is not allowed to be used"))
			}
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		ppp := pp.Index(i)
		for j, port := range co.Ports {
			allowed := false
			if port.HostPort == 0 {
				allowed = true
			}
			for _, r := range v.allowedHostPorts {
				if r.Min <= port.HostPort && port.HostPort <= r.Max {
					allowed = true
				}
			}
			if !allowed {
				errs = append(errs, field.Forbidden(ppp.Child("ports").Index(j), "Host port is not allowed to be used"))
			}
		}
	}
	return errs
}
