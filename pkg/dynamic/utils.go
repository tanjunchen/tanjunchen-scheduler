package dynamic

import (
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfigEnv (optionally) specify the location of kubeconfig file
const KubeConfigEnv = "KUBECONFIG"

// NewClusterConfig new kubeconfig
func NewClusterConfig(kubeconfig string) (*rest.Config, error) {
	var (
		cfg *rest.Config
		err error
	)

	if kubeconfig == "" {
		kubeconfig = os.Getenv(KubeConfigEnv)
	}

	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("error creating config from specified file: %v %w", kubeconfig, err)
		}
	} else {
		if cfg, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	}

	cfg.Timeout = time.Minute * 3
	cfg.QPS = 500
	cfg.Burst = 500

	return cfg, nil
}
