package hooks

import "github.com/cybozu-go/pod-security-admission/hooks/validators"

// SecurityProfile is a config for pod-security-admission
type SecurityProfile struct {
	Name string `json:"name"`

	HostNamespace        bool              `json:"hostNamespace"`
	PrivilegedContainers bool              `json:"privilegedContainers"`
	Capabilities         CapabilityProfile `json:"capabilities"`
	Volumes              VolumeProfile     `json:"volumes"`
	HostPorts            HostPortProfile   `json:"hostPorts"`
	UnsafeAppArmor       bool              `json:"unsafeAppArmor"`
	UnsafeSELinux        bool              `json:"unsafeSELinux"`
	UnsafeProcMount      bool              `json:"unsafeProcMount"`
	UnsafeSysctls        bool              `json:"unsafeSysctls"`
	PrivilegeEscalation  bool              `json:"privilegeEscalation"`
	Users                UserProfile       `json:"users"`
	RootGroups           bool              `json:"rootGroups"`
	UnsafeSeccomp        bool              `json:"unsafeSeccomp"`
}

type VolumeProfile struct {
	HostPathVolumes    bool `json:"hostPathVolumes"`
	NonCoreVolumeTypes bool `json:"nNonCoreVolumeTypes"`
}

type HostPortProfile struct {
	HostPorts        bool                   `json:"hostPorts"`
	AllowedHostPorts []validators.PortRange `json:"allowedHostPorts"`
}

type CapabilityProfile struct {
	UnsafeCapabilities  bool     `json:"unsafeCapabilities"`
	AllowedCapabilities []string `json:"allowedCapabilities"`
}

type UserProfile struct {
	RunAsRoot         bool `json:"runAsRoot"`
	ForceRunAsNonRoot bool `json:"forceRunAsNonRoot"`
}
