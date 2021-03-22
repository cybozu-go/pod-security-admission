package hooks

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = Describe("validate Pod webhook", func() {
	It("should allow a Pod having trusted images", func() {
		podManifest := `apiVersion: v1
kind: Pod
metadata:
  name: test-pod-validate-1
  namespace: default
spec:
  containers:
  - name: ubuntu
    image: quay.io/cybozu/ubuntu
  hostNetwork: true
`

		d := yaml.NewYAMLOrJSONDecoder(strings.NewReader(podManifest), 4096)
		po := &v1.Pod{}
		err := d.Decode(po)
		Expect(err).NotTo(HaveOccurred())

		err = k8sClient.Create(testCtx, po)
		Expect(err).NotTo(HaveOccurred())
	})
})
