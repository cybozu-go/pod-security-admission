package hooks

import "github.com/cybozu-go/pod-security-admission/hooks/validators"

// SecurityProfile is a config for pod-security-admission
type SecurityProfile struct {
	Name string `json:"name"`

	DenyHostNamespace        bool              `json:"denyHostNamespace"`
	DenyPrivilegedContainers bool              `json:"denyPrivilegedContainers"`
	Capabilities             CapabilityProfile `json:"capabilities"`
	Volumes                  VolumeProfile     `json:"volumes"`
	HostPorts                HostPortProfile   `json:"hostPorts"`
	DenyUnsafeAppArmor       bool              `json:"denyUnsafeAppArmor"`
	DenyUnsafeSELinux        bool              `json:"denyUnsafeSELinux"`
	DenyUnsafeProcMount      bool              `json:"denyUnsafeProcMount"`
	DenyUnsafeSysctls        bool              `json:"denyUnsafeSysctls"`
	DenyPrivilegeEscalation  bool              `json:"denyPrivilegeEscalation"`
	Users                    UserProfile       `json:"users"`
	DenyRootGroups           bool              `json:"denyRootGroups"`
	DenyUnsafeSeccomp        bool              `json:"denyUnsafeSeccomp"`
}

type VolumeProfile struct {
	DenyHostPathVolumes    bool `json:"denyHostPathVolumes"`
	DenyNonCoreVolumeTypes bool `json:"denyNonCoreVolumeTypes"`
}

type HostPortProfile struct {
	DenyHostPorts    bool                   `json:"denyHostPorts"`
	AllowedHostPorts []validators.PortRange `json:"allowedHostPorts"`
}

type CapabilityProfile struct {
	DenyUnsafeCapabilities bool     `json:"denyUnsafeCapabilities"`
	AllowedCapabilities    []string `json:"allowedCapabilities"`
}

type UserProfile struct {
	DenyRunAsRoot     bool `json:"denyRunAsRoot"`
	ForceRunAsNonRoot bool `json:"forceRunAsNonRoot"`
}
