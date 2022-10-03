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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	baselineMutatingWebhookPath     = "/mutate-baseline"
	baselineValidatingWebhookPath   = "/validate-baseline"
	hostpathMutatingWebhookPath     = "/mutate-hostpath"
	hostpathValidatingWebhookPath   = "/validate-hostpath"
	restrictedMutatingWebhookPath   = "/mutate-restricted"
	restrictedValidatingWebhookPath = "/validate-restricted"
	mutatingMutatingWebhookPath     = "/mutate-mutating"
	mutatingValidatingWebhookPath   = "/validate-mutating"
)

var k8sConfig *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testCtx = context.Background()
var testCancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Webhook Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	testCtx, testCancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{
				filepath.Join("testdata", "config"),
			},
		},
	}
	testEnv.ControlPlane.GetAPIServer().Configure().Append("feature-gates", "ProcMountType=true")

	var err error
	k8sConfig, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sConfig).NotTo(BeNil())

	scheme := runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = admissionv1beta1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = admissionv1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(k8sConfig, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(k8sConfig, ctrl.Options{
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
		Name: "baseline",
		AdditionalCapabilities: []string{
			"SYSLOG",
		},
		NonCoreVolumeTypes: true,
		AllowedHostPorts: []validators.PortRange{
			{
				Min: 65500,
				Max: 65502,
			},
		},
		RootGroups:               true,
		Seccomp:                  true,
		AllowPrivilegeEscalation: true,
		RunAsRoot:                true,
	}
	wh.Register(baselineValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(baselineValidatingWebhookPath), dec, baselineProfile))
	wh.Register(baselineMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(baselineMutatingWebhookPath), dec, baselineProfile))

	// "hostpath" profile = "baseline" profile + AllowedHostPaths
	hostpathProfile := SecurityProfile{
		Name: "hostpath",
		AdditionalCapabilities: []string{
			"SYSLOG",
		},
		NonCoreVolumeTypes: true,
		AllowedHostPaths: []validators.AllowedHostPath{
			{
				PathPrefix: "/etc/hos", // not "host"
			},
		},
		AllowedHostPorts: []validators.PortRange{
			{
				Min: 65500,
				Max: 65502,
			},
		},
		RootGroups:               true,
		Seccomp:                  true,
		AllowPrivilegeEscalation: true,
		RunAsRoot:                true,
	}
	wh.Register(hostpathValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(hostpathValidatingWebhookPath), dec, hostpathProfile))
	wh.Register(hostpathMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(hostpathMutatingWebhookPath), dec, hostpathProfile))

	restrictedProfile := SecurityProfile{
		Name: "restricted",
	}
	wh.Register(restrictedValidatingWebhookPath, NewPodValidator(mgr.GetClient(), ctrl.Log.WithName(restrictedValidatingWebhookPath), dec, restrictedProfile))
	wh.Register(restrictedMutatingWebhookPath, NewPodMutator(mgr.GetClient(), ctrl.Log.WithName(restrictedMutatingWebhookPath), dec, restrictedProfile))

	mutatingProfile := SecurityProfile{
		Name:              "mutating",
		ForceRunAsNonRoot: true,
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
})

var _ = AfterSuite(func() {
	testCancel()
	time.Sleep(10 * time.Millisecond)
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
