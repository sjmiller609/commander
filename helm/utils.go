package helm

import (
	"sync"

	"github.com/pkg/errors"
	"helm.sh/helm/pkg/chart"
	"helm.sh/helm/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	config     genericclioptions.RESTClientGetter
	configOnce sync.Once
)

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}


func GetNamespace() string {
	if ns, _, err := KubeConfig().ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

func KubeConfig() genericclioptions.RESTClientGetter {
	configOnce.Do(func() {
		config = kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
	})
	return config
}
