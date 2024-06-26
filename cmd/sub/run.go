package sub

import (
	"github.com/cybozu-go/pod-security-admission/hooks"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// +kubebuilder:scaffold:scheme
}

func run(addr string, port int, profs []hooks.SecurityProfile) error {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&config.zapOpts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: config.metricsAddr,
		},
		HealthProbeBindAddress: config.probeAddr,
		LeaderElection:         false,
		WebhookServer: webhook.NewServer(webhook.Options{
			Host:    addr,
			Port:    port,
			CertDir: config.certDir,
		}),
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	// register webhook handlers
	// admission.NewDecoder never returns non-nil error
	dec := admission.NewDecoder(scheme)

	wh := mgr.GetWebhookServer()
	for _, prof := range profs {
		wh.Register("/mutate-"+prof.Name, hooks.NewPodMutator(mgr.GetClient(), ctrl.Log.WithName("mutate-"+prof.Name), dec, prof))
		wh.Register("/validate-"+prof.Name, hooks.NewPodValidator(mgr.GetClient(), ctrl.Log.WithName("validate-"+prof.Name), dec, prof))
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		return err
	}
	return nil
}
