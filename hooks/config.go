package hooks

import "github.com/cybozu-go/pod-security-admission/hooks/validators"

// SecurityProfile is a config for pod-security-admission
type SecurityProfile struct {
	Name string `json:"name"`

	HostNamespace            bool                   `json:"hostNamespace"`
	Privileged               bool                   `json:"privileged"`
	Capabilities             bool                   `json:"capabilities"`
	AdditionalCapabilities   []string               `json:"additionalCapabilities"`
	HostPathVolumes          bool                   `json:"hostPathVolumes"`
	NonCoreVolumeTypes       bool                   `json:"nonCoreVolumeTypes"`
	HostPorts                bool                   `json:"hostPorts"`
	AllowedHostPorts         []validators.PortRange `json:"allowedHostPorts"`
	AppArmor                 bool                   `json:"appArmor"`
	SELinux                  bool                   `json:"seLinux"`
	ProcMount                bool                   `json:"procMount"`
	Sysctls                  bool                   `json:"sysctls"`
	AllowPrivilegeEscalation bool                   `json:"allowPrivilegeEscalation"`
	RunAsRoot                bool                   `json:"runAsRoot"`
	ForceRunAsNonRoot        bool                   `json:"forceRunAsNonRoot"`
	RootGroups               bool                   `json:"rootGroups"`
	Seccomp                  bool                   `json:"seccomp"`
}
