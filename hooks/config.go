package hooks

// SecurityProfile is a config for pod-security-admission
type SecurityProfile struct {
	Name       string   `json:"name"`
	Validators []string `json:"validators"`
	Mutators   []string `json:"mutators"`
}
