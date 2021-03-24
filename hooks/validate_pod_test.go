package hooks

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
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
			fmt.Printf("%s/%s: %s\n", po.Namespace, po.Name, statusErr.ErrStatus.Message)
		}
	}
}

var _ = Describe("validate Pod webhook", func() {
	It("should allow all pods in privileged namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "privileged", true)
		validatePod(filepath.Join("testdata", "baseline"), "privileged", true)
		validatePod(filepath.Join("testdata", "restricted"), "privileged", true)
	})

	It("should deny privileged pods in baseline namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "baseline", false)
		validatePod(filepath.Join("testdata", "baseline"), "baseline", true)
		validatePod(filepath.Join("testdata", "restricted"), "baseline", true)
	})

	It("should deny privileged and baseline pods in restricted namespace", func() {
		validatePod(filepath.Join("testdata", "privileged"), "restricted", false)
		validatePod(filepath.Join("testdata", "baseline"), "restricted", false)
		validatePod(filepath.Join("testdata", "restricted"), "restricted", true)
	})
})
