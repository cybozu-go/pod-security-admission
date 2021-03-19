package cmd

import (
	"flag"
	"fmt"
	"github.com/cybozu-go/pod-security-admission/hooks"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"
)

var config struct {
	metricsAddr           string
	probeAddr             string
	webhookAddr           string
	certDir               string
	configPath            string
	zapOpts               zap.Options
}

var rootCmd = &cobra.Command{
	Use:   "pod-security-admission",
	Short: "admission webhooks to ensure pod security standards",
	Long:  `Admission webhooks to ensure pod security standards.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		h, p, err := net.SplitHostPort(config.webhookAddr)
		if err != nil {
			return fmt.Errorf("invalid webhook address: %s, %v", config.webhookAddr, err)
		}
		numPort, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("invalid webhook address: %s, %v", config.webhookAddr, err)
		}
		conf, err := parseConfig(config.configPath)
		if err != nil {
			return err
		}
		return run(h, numPort, conf)
	},
}

func parseConfig(configPath string) (*hooks.Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var conf hooks.Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

// Execute executes the command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	fs := rootCmd.Flags()
	fs.StringVar(&config.metricsAddr, "metrics-addr", ":8080", "Listen address for metrics")
	fs.StringVar(&config.probeAddr, "health-probe-addr", ":8081", "Listen address for health probes")
	fs.StringVar(&config.webhookAddr, "webhook-addr", ":9443", "Listen address for the webhook endpoint")
	fs.StringVar(&config.certDir, "cert-dir", "", "certificate directory")
	fs.StringVar(&config.configPath, "config-path", "/etc/pod-security-admission/config.yaml", "Configuration for webhooks")

	goflags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(goflags)
	config.zapOpts.BindFlags(goflags)

	fs.AddGoFlagSet(goflags)
}
