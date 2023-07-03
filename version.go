package psa

var version = "unset"

// Version returns the semantic versioning string of pod-security-admission.
func Version() string {
	return version
}
