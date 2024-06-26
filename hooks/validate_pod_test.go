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

func validatePod(dir string, namespace string, allowed bool) {
	entries, err := os.ReadDir(dir)
	Expect(err).NotTo(HaveOccurred())
	for _, e := range entries {
		y, err := os.ReadFile(filepath.Join(dir, e.Name()))
		Expect(err).NotTo(HaveOccurred())
		d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(y), 4096)
		po := &corev1.Pod{}
		err = d.Decode(po)
		Expect(err).NotTo(HaveOccurred())
		po.Namespace = namespace
		err = k8sClient.Create(testCtx, po)
		if allowed {
			Expect(err).NotTo(HaveOccurred(), "pod: %v", po)
		} else {
			Expect(err).To(HaveOccurred(), "pod: %v", po)
			statusErr := &k8serrors.StatusError{}
			Expect(errors.As(err, &statusErr)).To(BeTrue())
			expected := po.Annotations["test.pod-security.cybozu.com/message"]
			Expect(statusErr.ErrStatus.Message).To(ContainSubstring(expected))
		}
	}
}

var _ = Describe("validate Pod webhook", func() {
	It("should allow all pods in privileged namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "privileged", true)
		validatePod(filepath.Join("testdata", "hostpath"), "privileged", true)
		validatePod(filepath.Join("testdata", "baseline"), "privileged", true)
		validatePod(filepath.Join("testdata", "restricted"), "privileged", true)
	})
	It("should deny privileged pods in hostpath namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "hostpath", false)
		validatePod(filepath.Join("testdata", "hostpath"), "hostpath", true)
		validatePod(filepath.Join("testdata", "baseline"), "hostpath", true)
		validatePod(filepath.Join("testdata", "restricted"), "hostpath", true)
	})
	It("should deny privileged and hostpath pods in baseline namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "baseline", false)
		validatePod(filepath.Join("testdata", "hostpath"), "baseline", false)
		validatePod(filepath.Join("testdata", "baseline"), "baseline", true)
		validatePod(filepath.Join("testdata", "restricted"), "baseline", true)
	})
	It("should deny privileged, hostpath, and baseline pods in restricted namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "restricted", false)
		validatePod(filepath.Join("testdata", "hostpath"), "restricted", false)
		validatePod(filepath.Join("testdata", "baseline"), "restricted", false)
		validatePod(filepath.Join("testdata", "restricted"), "restricted", true)
	})
})
