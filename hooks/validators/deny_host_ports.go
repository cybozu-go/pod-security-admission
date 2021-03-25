package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyHostPorts is a Validator that denies usage of HostPorts
func DenyHostPorts(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec")
	var errs field.ErrorList

	pp := p.Child("containers")
	for i, co := range pod.Spec.Containers {
		ppp := pp.Index(i)
		for j, port := range co.Ports {
			if port.HostPort != 0 {
				errs = append(errs, field.Forbidden(ppp.Index(j), "Host port is not allowed to be used"))
			}
		}
	}

	pp = p.Child("initContainers")
	for i, co := range pod.Spec.InitContainers {
		ppp := pp.Index(i)
		for j, port := range co.Ports {
			if port.HostPort != 0 {
				errs = append(errs, field.Forbidden(ppp.Index(j), "Host port is not allowed to be used"))
			}
		}
	}
	return errs
}
