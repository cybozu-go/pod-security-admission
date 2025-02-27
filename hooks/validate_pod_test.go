package hooks

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func validatePod(dir string, profile string) {
	entries, err := os.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred())
	for _, e := range entries {
		By("validating " + e.Name() + " in " + profile)
		y, err := os.ReadFile(filepath.Join(dir, e.Name()))
		Expect(err).NotTo(HaveOccurred())
		d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(y), 4096)
		po := &corev1.Pod{}
		err = d.Decode(po)
		Expect(err).NotTo(HaveOccurred())
		po.Namespace = profile
		err = k8sClient.Create(testCtx, po)
		expected := po.Annotations["expected.pod-security.cybozu.com/"+profile]
		if expected == "" {
			// If `expected` is empty, it means that the pod is allowed by the profile.
			Expect(err).NotTo(HaveOccurred(), "pod: %v", po)
		} else {
			// If `expected` is not empty, it means that the pod is denied by the profile.
			// `expected` contains an error message.
			Expect(err).To(HaveOccurred(), "pod: %v", po)
			statusErr := &k8serrors.StatusError{}
			Expect(errors.As(err, &statusErr)).To(BeTrue())
			Expect(statusErr.ErrStatus.Message).To(ContainSubstring(expected))
		}
	}
}

var _ = Describe("validate Pod webhook", func() {
	for _, profile := range []string{"privileged", "hostpath", "baseline", "restricted"} {
		It("should validate pods in "+profile, func() {
			validatePod(filepath.Join("testdata", "validating"), profile)
		})
	}

})
