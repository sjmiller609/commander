package kubernetes

import (
	"helm.sh/helm/pkg/cli"
	"helm.sh/helm/pkg/kube"

	//"helm.sh/helm/pkg/kube"
	//"k8s.io/cli-runtime/pkg/genericclioptions"

	"sync"

	// Required to authenticate against GKE clusters
	//_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// Import to initialize client auth plugins.
)

var (
	settings   cli.EnvSettings
	//config     genericclioptions.RESTClientGetter
	configOnce sync.Once
)

//func GetKubeConfig() (*rest.Config, error) {
//	logger := log.WithField("function", "GetKubeConfig")
//	logger.Debug("Getting KubeConfig")
//
//	var configErr error = nil
//
//	if appConfig.KubeConfig != "" {
//		logger.Debug("Using config at ", appConfig.KubeConfig)
//		kubeConfig, configErr = clientcmd.BuildConfigFromFlags("", appConfig.KubeConfig)
//	} else {
//		logger.Debug("Using in-cluster config")
//		kubeConfig, configErr = rest.InClusterConfig()
//	}
//	return kubeConfig, configErr
//}

func GetKubeConfig() {
	configOnce.Do(func() {
		kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
	})
}