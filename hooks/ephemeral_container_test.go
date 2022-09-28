package hooks

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

var _ = Describe("Admission Webhook for EphemeralContainer", func() {

	unmasked := corev1.UnmaskedProcMount
	DescribeTable("Validating Ephemeral Container", func(namespace, name string, ec corev1.EphemeralContainer, allowed bool, message string) {
		podManifest := `apiVersion: v1
kind: Pod
metadata:
  namespace: %s
  name: %s
spec:
  containers:
  - name: ubuntu
    image: quay.io/cybozu/ubuntu
    securityContext:
      runAsNonRoot: true
`

		pod := &corev1.Pod{}
		d := yaml.NewYAMLOrJSONDecoder(strings.NewReader(fmt.Sprintf(podManifest, namespace, name)), 4096)
		err := d.Decode(pod)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Create(testCtx, pod)
		Expect(err).NotTo(HaveOccurred())

		k8s, err := kubernetes.NewForConfig(k8sConfig)
		Expect(err).NotTo(HaveOccurred())
		podClient := k8s.CoreV1().Pods(pod.Namespace)
		pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, ec)

		_, err = podClient.UpdateEphemeralContainers(testCtx, pod.Name, pod, metav1.UpdateOptions{})
		if allowed {
			Expect(err).NotTo(HaveOccurred(), "pod: %v", pod)
		} else {
			Expect(err).To(HaveOccurred(), "pod: %v", pod)
			statusErr := &k8serrors.StatusError{}
			Expect(errors.As(err, &statusErr)).To(BeTrue())
			Expect(statusErr.ErrStatus.Message).To(ContainSubstring(message))
		}
	},
		Entry("Valid Ephemeral Container", "restricted", "test-simple-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
				},
			},
		}, true, ""),
		Entry("Privileged Ephemeral Container", "baseline", "test-privileged-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					Privileged:   pointer.Bool(true),
				},
			},
		}, false, "denied the request: spec.ephemeralContainer[0].securityContext.privileged: Forbidden: Privileged containers are not allowed"),
		Entry("AllowPrivilegeEscalation Ephemeral Container", "restricted", "test-allow-privilege-escalation-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot:             pointer.Bool(true),
					AllowPrivilegeEscalation: pointer.Bool(true),
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.allowPrivilegeEscalation: Forbidden: Allowing privilege escalation for containers is not allowed"),
		Entry("RootGroup Ephemeral Container", "restricted", "test-root-group-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					RunAsGroup:   pointer.Int64(0),
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.runAsGroup: Forbidden: Running with the root GID is forbidden"),
		// runAsNonRoot of an ephemeral container will not be validated until the following issue is completed.
		// https://github.com/kubernetes/kubectl/issues/1108
		/*
			Entry("RunAsRoot Ephemeral Container", "restricted", "test-run-as-root-ec", corev1.EphemeralContainer{
				EphemeralContainerCommon: corev1.EphemeralContainerCommon{
					Name:  "debug",
					Image: "quay.io/cybozu/ubuntu-debug",
					SecurityContext: &corev1.SecurityContext{
						RunAsNonRoot: pointer.Bool(false),
					},
				},
			}, false, "denied the request: [spec.ephemeralContainers[0].securityContext.runAsNonRoot: Forbidden: RunAsNonRoot must be true, spec.securityContext: Forbidden: RunAsNonRoot must be true]"),
		*/
		Entry("UnsafeCapability Ephemeral Container", "restricted", "test-unsafe-capability-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					Capabilities: &corev1.Capabilities{
						Add: []corev1.Capability{
							"SYSLOG",
						},
					},
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.capabilities.add[0]: Forbidden: Adding capability SYSLOG is not allowed"),
		Entry("UnsafeProcMount Ephemeral Container", "restricted", "test-unsafe-procmount-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					ProcMount:    &unmasked,
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.procMount: Forbidden: ProcMountType Unmasked is not allowed"),
		Entry("UnsafeSeccomp Ephemeral Container", "restricted", "test-unsafe-seccomp-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					SeccompProfile: &corev1.SeccompProfile{
						Type:             corev1.SeccompProfileTypeLocalhost,
						LocalhostProfile: pointer.String("profiles/audit.json"),
					},
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.seccompProfile.type: Forbidden: Localhost is not an allowed seccomp profile"),
		Entry("UnsafeSELinux Ephemeral Container", "restricted", "test-unsafe-selinux-ec", corev1.EphemeralContainer{
			EphemeralContainerCommon: corev1.EphemeralContainerCommon{
				Name:  "debug",
				Image: "quay.io/cybozu/ubuntu-debug",
				SecurityContext: &corev1.SecurityContext{
					RunAsNonRoot: pointer.Bool(true),
					SELinuxOptions: &corev1.SELinuxOptions{
						Level: "s0:c123,c456",
					},
				},
			},
		}, false, "denied the request: spec.ephemeralContainers[0].securityContext.selinuxOptions: Forbidden: Setting custom SELinux options is not allowed"),
	)

	// runAsNonRoot of an ephemeral container will not be mutated until the following issue is completed.
	// https://github.com/kubernetes/kubectl/issues/1108
	/*
			Context("Mutating EphemeralContainer", func() {
				It("", func() {
					podManifest := `apiVersion: v1
		kind: Pod
		metadata:
		  namespace: %s
		  name: %s
		spec:
		  containers:
		  - name: ubuntu
		    image: quay.io/cybozu/ubuntu
		    securityContext:
		      runAsNonRoot: true
		`

					pod := &corev1.Pod{}
					d := yaml.NewYAMLOrJSONDecoder(strings.NewReader(fmt.Sprintf(podManifest, "mutating", "test-mutate-ec")), 4096)
					err := d.Decode(pod)
					Expect(err).NotTo(HaveOccurred())
					err = k8sClient.Create(testCtx, pod)
					Expect(err).NotTo(HaveOccurred())

					k8s, err := kubernetes.NewForConfig(k8sConfig)
					Expect(err).NotTo(HaveOccurred())
					podClient := k8s.CoreV1().Pods(pod.Namespace)
					ec := corev1.EphemeralContainer{
						EphemeralContainerCommon: corev1.EphemeralContainerCommon{
							Name:  "debug",
							Image: "quay.io/cybozu/ubuntu-debug",
						},
					}
					pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, ec)

					_, err = podClient.UpdateEphemeralContainers(testCtx, pod.Name, pod, metav1.UpdateOptions{})
					Expect(err).NotTo(HaveOccurred(), "pod: %v", pod)

					ret := &corev1.Pod{}
					err = k8sClient.Get(testCtx, types.NamespacedName{Name: pod.GetName(), Namespace: pod.GetNamespace()}, ret)
					Expect(err).NotTo(HaveOccurred())
					Expect(pod.Spec.EphemeralContainers).Should(HaveLen(1))
					Expect(pod.Spec.EphemeralContainers[0].SecurityContext).ShouldNot(BeNil())
					Expect(pod.Spec.EphemeralContainers[0].SecurityContext.RunAsNonRoot).ShouldNot(BeNil())
					Expect(pod.Spec.EphemeralContainers[0].SecurityContext.RunAsNonRoot).ShouldNot(HaveValue(Equal(true)))
				})
			})
	*/
})
