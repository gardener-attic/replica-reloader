// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package kubeconfig

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeConfig struct {
	kubeconfig string
}

// AddFlags adds flags related to kubeconfig configuration to the specified FlagSet.
func (k *KubeConfig) AddFlags(fs *pflag.FlagSet) {
	if k == nil {
		return
	}

	fs.StringVar(&k.kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster")
}

// Complete creates a *rest.Config for talking to a Kubernetes API server.
func (k *KubeConfig) Complete() (*rest.Config, error) {
	if k == nil {
		return nil, fmt.Errorf("KubeConfig cannot be nil")
	}

	if len(k.kubeconfig) > 0 {
		return loadWithOverride(&clientcmd.ClientConfigLoadingRules{ExplicitPath: k.kubeconfig})
	}

	kubeconfigPath := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(kubeconfigPath) == 0 {
		if c, err := rest.InClusterConfig(); err == nil {
			return c, nil
		}
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if _, ok := os.LookupEnv("HOME"); !ok {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("could not get current user: %v", err)
		}
		loadingRules.Precedence = append(loadingRules.Precedence, path.Join(u.HomeDir, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
	}

	cfg, err := loadWithOverride(loadingRules)
	if err != nil {
		return nil, err
	}

	cfg.QPS = 2.0
	cfg.Burst = 4.0

	return cfg, nil
}

func loadWithOverride(loader clientcmd.ClientConfigLoader) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &clientcmd.ConfigOverrides{}).ClientConfig()
}
