package helm

import (
	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/storage"
	"helm.sh/helm/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

func NewActionConfig(kubeConfig genericclioptions.RESTClientGetter, kubeClient genericclioptions.RESTClientGetter) *action.Configuration {
	var store *storage.Storage
	switch os.Getenv("HELM_DRIVER") {
	case "memory":
		d := driver.NewMemory()
		store = storage.Init(d)
	default:
		// Not sure what to do here.
		panic("Unknown driver in HELM_DRIVER: " + os.Getenv("HELM_DRIVER"))
	}

	return &action.Configuration{
		RESTClientGetter: kubeConfig,
		KubeClient:       kubeClient,
		Releases:         store,
	}
}