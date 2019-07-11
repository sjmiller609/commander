package kubernetes

import (
	"helm.sh/helm/pkg/cli"
	"helm.sh/helm/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"sync"

	// Required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// Import to initialize client auth plugins.
)

var (
	settings   cli.EnvSettings
	kubeConfig     genericclioptions.RESTClientGetter
	configOnce sync.Once
)

func GetKubeConfig() genericclioptions.RESTClientGetter {
	configOnce.Do(func() {
		kubeConfig = kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
	})
	return kubeConfig
}