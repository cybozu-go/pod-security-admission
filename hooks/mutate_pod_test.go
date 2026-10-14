package hooks

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"
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

		ret := &corev1.Pod{}
		err = k8sClient.Get(testCtx, types.NamespacedName{Name: po.GetName(), Namespace: po.GetNamespace()}, ret)
		Expect(err).NotTo(HaveOccurred())
	}
}

var _ = Describe("mutate Pod webhook", func() {
	It("should mutate pods in mutating namespace", func() {
		mutatePod(filepath.Join("testdata", "mutating"), "mutating")
	})

	It("should not force pods that use user namespaces to run as non-root", func() {
		newPod := func(name string, hostUsers *bool) *corev1.Pod {
			return &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "mutating",
					Name:      name,
				},
				Spec: corev1.PodSpec{
					HostUsers: hostUsers,
					Containers: []corev1.Container{
						{Name: "ubuntu", Image: "ghcr.io/cybozu/ubuntu"},
					},
					InitContainers: []corev1.Container{
						{Name: "init", Image: "ghcr.io/cybozu/ubuntu"},
					},
				},
			}
		}

		By("creating a pod that uses user namespaces")
		po := newPod("userns-not-mutated", ptr.To(false))
		Expect(k8sClient.Create(testCtx, po)).To(Succeed())
		ret := &corev1.Pod{}
		Expect(k8sClient.Get(testCtx, types.NamespacedName{Name: po.Name, Namespace: po.Namespace}, ret)).To(Succeed())
		for _, co := range append(ret.Spec.Containers, ret.Spec.InitContainers...) {
			if co.SecurityContext != nil {
				Expect(co.SecurityContext.RunAsNonRoot).To(BeNil(), "container: %s", co.Name)
			}
		}

		By("creating a pod that does not use user namespaces")
		po = newPod("no-userns-mutated", nil)
		Expect(k8sClient.Create(testCtx, po)).To(Succeed())
		ret = &corev1.Pod{}
		Expect(k8sClient.Get(testCtx, types.NamespacedName{Name: po.Name, Namespace: po.Namespace}, ret)).To(Succeed())
		for _, co := range append(ret.Spec.Containers, ret.Spec.InitContainers...) {
			Expect(co.SecurityContext).NotTo(BeNil(), "container: %s", co.Name)
			Expect(co.SecurityContext.RunAsNonRoot).To(HaveValue(BeTrue()), "container: %s", co.Name)
		}
	})
})
