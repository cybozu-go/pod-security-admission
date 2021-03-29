package hooks

import "github.com/cybozu-go/pod-security-admission/hooks/validators"

// SecurityProfile is a config for pod-security-admission
type SecurityProfile struct {
	Name string `json:"name"`

	DenyHostNamespace        bool `json:"denyHostNamespace"`
	DenyPrivilegedContainers bool `json:"denyPrivilegedContainers"`
	DenyUnsafeCapabilities   bool `json:"denyUnsafeCapabilities"`
	DenyHostPathVolumes      bool `json:"denyHostPathVolumes"`
	DenyHostPorts            bool `json:"denyHostPorts"`
	DenyUnsafeAppArmor       bool `json:"denyUnsafeAppArmor"`
	DenyUnsafeSELinux        bool `json:"denyUnsafeSELinux"`
	DenyUnsafeProcMount      bool `json:"denyUnsafeProcMount"`
	DenyUnsafeSysctls        bool `json:"denyUnsafeSysctls"`
	DenyNonCoreVolumeTypes   bool `json:"denyNonCoreVolumeTypes"`
	DenyPrivilegeEscalation  bool `json:"denyPrivilegeEscalation"`
	DenyRunAsRoot            bool `json:"denyRunAsRoot"`
	DenyRootGroups           bool `json:"denyRootGroups"`
	DenyUnsafeSeccomp        bool `json:"denyUnsafeSeccomp"`

	AllowedCapabilities []string               `json:"allowedCapabilities"`
	AllowedHostPorts    []validators.PortRange `json:"allowedHostPorts"`

	ForceRunAsNonRoot bool `json:"forceRunAsNonRoot"`
}
