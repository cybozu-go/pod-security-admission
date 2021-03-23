package hooks

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podValidator struct {
	client     client.Client
	decoder    *admission.Decoder
	validators []string
}

// NewPodValidator creates a webhook handler for Pod.
func NewPodValidator(c client.Client, dec *admission.Decoder, validators []string) http.Handler {
	v := &podValidator{
		client:     c,
		decoder:    dec,
		validators: validators,
	}
	return &webhook.Admission{Handler: v}
}

func (v *podValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	po := &corev1.Pod{}
	err := v.decoder.Decode(req, po)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for _, name := range v.validators {
		validator := validators[name]
		res := validator(po)
		if !res.Allowed {
			return res
		}
	}

	return admission.Allowed("ok")
}

type Validator func(pod *corev1.Pod) admission.Response

var validators = map[string]Validator{
	"deny-host-namespace":        denyHostNamespace,
	"deny-privileged-containers": denyPrivilegedContainers,
	"deny-unsafe-capabilities":   denyUnsafeCapabilities,
	"deny-host-path-volumes":     denyHostPathVolumes,
	"deny-host-ports":            denyHostPorts,
	"deny-unsafe-apparmor":       denyUnsafeAppArmor,
	"deny-unsafe-selinux":        denyUnsafeSELinux,
	"deny-unsafe-proc-mount":     denyUnsafeProcMount,
	"deny-unsafe-sysctls":        denyUnsafeSysctls,
	"deny-non-core-volume-types": denyNonCoreVolumeTypes,
	"deny-privilege-escalation":  denyPrivilegeEscalation,
	"deny-run-as-root":           denyRunAsRoot,
	"deny-root-groups":           denyRootGroups,
	"deny-unsafe-seccomp":        denyUnsafeSeccomp,
}

func denyHostNamespace(pod *corev1.Pod) admission.Response {
	if pod.Spec.HostNetwork != false {
		return admission.Denied("Host network is not allowed to be used")
	}
	if pod.Spec.HostPID != false {
		return admission.Denied("Host pid is not allowed to be used")
	}
	if pod.Spec.HostIPC != false {
		return admission.Denied("Host ipc is not allowed to be used")
	}
	return admission.Allowed("ok")
}

func denyPrivilegedContainers(pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, c := range containers {
		if c.SecurityContext == nil || c.SecurityContext.Privileged == nil {
			continue
		}
		if *c.SecurityContext.Privileged == true {
			return admission.Denied("Privileged containers are not allowed")
		}
	}
	return admission.Allowed("ok")
}

var defaultCapabilities = []string{
	"AUDIT_WRITE",
	"CHOWN",
	"DAC_OVERRIDE",
	"FOWNER",
	"FSETID",
	"KILL",
	"MKNOD",
	"NET_BIND_SERVICE",
	"NET_RAW",
	"SETFCAP",
	"SETGID",
	"SETPCAP",
	"SETUID",
	"SYS_CHROOT",
}

func containsString(list []string, item string) bool {
	for _, c := range list {
		if item == c {
			return true
		}
	}
	return false
}

func denyUnsafeCapabilities(pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, c := range containers {
		if c.SecurityContext == nil || c.SecurityContext.Capabilities == nil {
			continue
		}
		for _, add := range c.SecurityContext.Capabilities.Add {
			if !containsString(defaultCapabilities, string(add)) {
				return admission.Denied(fmt.Sprintf("Adding capability %s is not allowed", add))
			}
		}
	}

	return admission.Allowed("ok")
}

func denyHostPathVolumes(pod *corev1.Pod) admission.Response {
	for _, vol := range pod.Spec.Volumes {
		if vol.HostPath == nil {
			continue
		}
		if len(vol.HostPath.Path) != 0 || vol.HostPath.Type != nil {
			return admission.Denied("Host path is not allowed to be used")
		}
	}
	return admission.Allowed("ok")
}

func denyHostPorts(pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.HostPort != 0 {
				return admission.Denied("Host port is not allowed to be used")
			}
		}
	}
	return admission.Allowed("ok")
}

func denyUnsafeAppArmor(pod *corev1.Pod) admission.Response {
	for k, v := range pod.Annotations {
		if strings.HasPrefix(k, corev1.AppArmorBetaContainerAnnotationKeyPrefix) &&
			v != corev1.AppArmorBetaProfileRuntimeDefault {
			return admission.Denied(fmt.Sprintf("%s is not an allowed AppArmor profile", v))
		}
	}
	return admission.Allowed("ok")
}

func denyUnsafeSELinux(pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SELinuxOptions != nil {
		return admission.Denied("Setting custom SELinux options is not allowed")
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext != nil && container.SecurityContext.SELinuxOptions != nil {
			return admission.Denied("Setting custom SELinux options is not allowed")
		}
	}
	return admission.Allowed("ok")
}

func denyUnsafeProcMount(pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.ProcMount == nil {
			continue
		}
		proc := *container.SecurityContext.ProcMount
		if proc != corev1.DefaultProcMount {
			return admission.Denied(fmt.Sprintf("ProcMountType %s is not allowed", proc))
		}
	}
	return admission.Allowed("ok")
}

var allowedSysctls = []string{
	"kernel.shm_rmid_forced",
	"net.ipv4.ip_local_port_range",
	"net.ipv4.tcp_syncookies",
	"net.ipv4.ping_group_range",
}

func denyUnsafeSysctls(pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext == nil {
		return admission.Allowed("ok")
	}
	for _, sysctl := range pod.Spec.SecurityContext.Sysctls {
		if !containsString(allowedSysctls, sysctl.Name) {
			return admission.Denied(fmt.Sprintf("Setting sysctl %s is not allowed", sysctl.Name))
		}
	}
	return admission.Allowed("ok")
}

func denyNonCoreVolumeTypes(pod *corev1.Pod) admission.Response {
	for _, vol := range pod.Spec.Volumes {
		if vol.HostPath != nil ||
			vol.GCEPersistentDisk != nil ||
			vol.AWSElasticBlockStore != nil ||
			vol.GitRepo != nil ||
			vol.NFS != nil ||
			vol.ISCSI != nil ||
			vol.Glusterfs != nil ||
			vol.RBD != nil ||
			vol.FlexVolume != nil ||
			vol.Cinder != nil ||
			vol.CephFS != nil ||
			vol.Flocker != nil ||
			vol.FC != nil ||
			vol.AzureFile != nil ||
			vol.VsphereVolume != nil ||
			vol.Quobyte != nil ||
			vol.AzureDisk != nil ||
			vol.PortworxVolume != nil ||
			vol.ScaleIO != nil ||
			vol.StorageOS != nil ||
			vol.CSI != nil {
			return admission.Denied(fmt.Sprintf("Volume type %s is not allowed to be used", vol.String()))
		}
	}
	return admission.Allowed("ok")
}

func denyPrivilegeEscalation(pod *corev1.Pod) admission.Response {
	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.AllowPrivilegeEscalation == nil {
			continue
		}
		escalation := *container.SecurityContext.AllowPrivilegeEscalation
		if escalation {
			return admission.Denied("Allowing privilege escalation for containers is not allowed")
		}
	}
	return admission.Allowed("ok")
}

func denyRunAsRoot(pod *corev1.Pod) admission.Response {
	validate := func(runAsNonRoot *bool, runAsUser *int64) admission.Response {
		if runAsNonRoot == nil && runAsUser == nil {
			return admission.Denied("RunAsNonRoot must be true")
		}
		if runAsNonRoot != nil && *runAsNonRoot == false {
			return admission.Denied("RunAsNonRoot must be true")
		}
		if runAsUser != nil && *runAsUser == 0 {
			return admission.Denied("Running with the root UID is forbidden")
		}
		return admission.Allowed("ok")
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	allContainersAllowed := true
	for _, container := range containers {
		if container.SecurityContext == nil || (container.SecurityContext.RunAsNonRoot == nil && container.SecurityContext.RunAsUser == nil) {
			allContainersAllowed = false
			continue
		}
		if res := validate(container.SecurityContext.RunAsNonRoot, container.SecurityContext.RunAsUser); !res.Allowed {
			return res
		}
	}
	if allContainersAllowed {
		return admission.Allowed("ok")
	}

	if pod.Spec.SecurityContext == nil {
		return admission.Denied("RunAsNonRoot must be true")
	}
	return validate(pod.Spec.SecurityContext.RunAsNonRoot, pod.Spec.SecurityContext.RunAsUser)
}

func denyRootGroups(pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil {
		if pod.Spec.SecurityContext.RunAsGroup != nil && *pod.Spec.SecurityContext.RunAsGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
		for _, group := range pod.Spec.SecurityContext.SupplementalGroups {
			if group == 0 {
				return admission.Denied("Running with the supplementary GID is forbidden")
			}
		}
		if pod.Spec.SecurityContext.FSGroup != nil && *pod.Spec.SecurityContext.FSGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil {
			continue
		}
		if container.SecurityContext.RunAsGroup != nil && *container.SecurityContext.RunAsGroup == 0 {
			return admission.Denied("Running with the root GID is forbidden")
		}
	}
	return admission.Allowed("ok")
}

func denyUnsafeSeccomp(pod *corev1.Pod) admission.Response {
	if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.SeccompProfile != nil && pod.Spec.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		return admission.Denied(fmt.Sprintf("%s is not an allowed seccomp prifile", pod.Spec.SecurityContext.SeccompProfile.Type))
	}

	containers := make([]corev1.Container, len(pod.Spec.Containers)+len(pod.Spec.InitContainers))
	copy(containers, pod.Spec.Containers)
	copy(containers[len(pod.Spec.Containers):], pod.Spec.InitContainers)

	for _, container := range containers {
		if container.SecurityContext == nil || container.SecurityContext.SeccompProfile == nil {
			continue
		}
		if container.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			return admission.Denied(fmt.Sprintf("%s is not an allowed seccomp prifile", container.SecurityContext.SeccompProfile.Type))
		}
	}
	return admission.Allowed("ok")
}
