package kubernetes

import (
	"k8s.io/client-go/rest"
	// Required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig *rest.Config = nil
)

func GetKubeConfig() (*rest.Config, error) {
	logger := log.WithField("function", "GetKubeConfig")
	logger.Debug("Getting KubeConfig")

	var configErr error = nil

	if appConfig.KubeConfig != "" {
		logger.Debug("Using config at ", appConfig.KubeConfig)
		kubeConfig, configErr = clientcmd.BuildConfigFromFlags("", appConfig.KubeConfig)
	} else {
		logger.Debug("Using in-cluster config")
		kubeConfig, configErr = rest.InClusterConfig()
	}
	return kubeConfig, configErr
}