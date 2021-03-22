package hooks

// Config is a config for pod-security-admission
type Config struct {
}

type SecurityProfile struct {
	Name       string   `json:"name"`
	Validators []string `json:"validators"`
	Mutators   []string `json:"mutators"`
}
