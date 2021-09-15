module github.com/cybozu-go/pod-security-admission

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/spf13/cobra v1.2.1
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20210802155522-efc7438f0176
	sigs.k8s.io/controller-runtime v0.9.7
	sigs.k8s.io/yaml v1.2.0
)
