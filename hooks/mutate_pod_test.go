package hooks

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func mutatePod(dir string, namespace string) {
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
		Expect(err).NotTo(HaveOccurred())

		ret := &v1.Pod{}
		err = k8sClient.Get(testCtx, types.NamespacedName{Name: po.GetName(), Namespace: po.GetNamespace()}, ret)
		Expect(err).NotTo(HaveOccurred())
	}
}

var _ = Describe("mutate Pod webhook", func() {
	It("should mutate pods in mutating namespace", func() {
		mutatePod(filepath.Join("testdata", "mutating"), "mutating")
	})
})
