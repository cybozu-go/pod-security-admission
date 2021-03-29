package hooks

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	//+kubebuilder:scaffold:imports
	"github.com/cybozu-go/pod-security-admission/hooks/validators"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	baselineMutatingWebhookPath     = "/mutate-baseline"
	baselineValidatingWebhookPath   = "/validate-baseline"
	restrictedMutatingWebhookPath   = "/mutate-restricted"
	restrictedValidatingWebhookPath = "/validate-restricted"
	mutatingMutatingWebhookPath     = "/mutate-mutating"
	mutatingValidatingWebhookPath   = "/validate-mutating"
)

var k8sClient client.Client
var testEnv *envtest.Environment
var testCtx = context.Background()
var testCancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Webhook Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	testCtx, testCancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{
				filepath.Join("..", "config", "webhook"),
				filepath.Join("testdata", "config"),
			},
		},
		KubeAPIServerFlags: append(envtest.DefaultKubeAPIServerFlags, "--feature-gates=ProcMountType=true"),
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = admissionv1beta1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = admissionv1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme,
		Host:               webhookInstallOptions.LocalServingHost,
		Port:               webhookInstallOptions.LocalServingPort,
		CertDir:            webhookInstallOptions.LocalServingCertDir,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	Expect(err).NotTo(HaveOccurred())

	dec, err := admission.NewDecoder(scheme)
	Expect(err).NotTo(HaveOccurred())
	wh := mgr.GetWebhookServer()
	baselineProfile := SecurityProfile{
		Name:                     "baseline",
		DenyHostNamespace:        true,
		DenyPrivilegedContainers: true,
		Capabilities: CapabilityProfile{
			DenyUnsafeCapabilities: true,
			AllowedCapabilities: []string{
				"SYSLOG",
			},
		},
		Volumes: VolumeProfile{DenyHostPathVolumes: true},
		HostPorts: HostPortProfile{
			DenyHostPorts: true,
			AllowedHostPorts: []validators.PortRange{
				{
					Min: 65500,
					Max: 65502,
				},
			},
		},
		DenyUnsafeAppArmor:  true,
		DenyUnsafeSELinux:   true,
		DenyUnsafeProcMount: true,
		DenyUnsafeSysctls:   true,
	}
	wh.Register(baselineValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(baselineValidatingWebhookPath), dec, baselineProfile))
	wh.Register(baselineMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(baselineMutatingWebhookPath), dec, baselineProfile))
	restrictedProfile := SecurityProfile{
		Name:                    "restricted",
		Volumes:                 VolumeProfile{DenyNonCoreVolumeTypes: true},
		DenyPrivilegeEscalation: true,
		RunAsRoot:               RunAsRootProfile{DenyRunAsRoot: true},
		DenyRootGroups:          true,
		DenyUnsafeSeccomp:       true,
	}
	wh.Register(restrictedValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(restrictedValidatingWebhookPath), dec, restrictedProfile))
	wh.Register(restrictedMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(restrictedMutatingWebhookPath), dec, restrictedProfile))

	mutatingProfile := SecurityProfile{
		Name: "mutating",
		RunAsRoot: RunAsRootProfile{
			DenyRunAsRoot:     true,
			ForceRunAsNonRoot: true,
		},
	}
	wh.Register(mutatingValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(mutatingValidatingWebhookPath), dec, mutatingProfile))
	wh.Register(mutatingMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(mutatingMutatingWebhookPath), dec, mutatingProfile))

	//+kubebuilder:scaffold:webhook

	go func() {
		err = mgr.Start(testCtx)
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}
	}()

	// wait for the webhook server to get ready
	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	Eventually(func() error {
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		conn.Close()
		return nil
	}).Should(Succeed())

	// create namespaces
	entries, err := os.ReadDir(filepath.Join("testdata", "namespace"))
	Expect(err).NotTo(HaveOccurred())
	for _, e := range entries {
		y, err := os.ReadFile(filepath.Join("testdata", "namespace", e.Name()))
		Expect(err).NotTo(HaveOccurred())
		d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(y), 4096)
		ns := &corev1.Namespace{}
		err = d.Decode(ns)
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.Create(testCtx, ns)
		Expect(err).NotTo(HaveOccurred())
	}
}, 60)

var _ = AfterSuite(func() {
	testCancel()
	time.Sleep(10 * time.Millisecond)
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
