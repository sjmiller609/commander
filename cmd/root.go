package cmd

import (
	//"helm.sh/helm/pkg/kube"
	"os"

	"github.com/astronomer/commander/helm"
	//"github.com/astronomer/commander/api"
	//"github.com/astronomer/commander/config"
	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/kube"
	"helm.sh/helm/pkg/storage"
	"helm.sh/helm/pkg/storage/driver"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log       = logrus.WithField("package", "cmd")
	//_ = kubernetes.KubeProvisioner{}
)

// RootCmd is the commander root command.
var RootCmd = &cobra.Command{
	Use: "commander",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	// Set up logging
	logrus.SetOutput(os.Stdout)
	logger := log.WithField("function", "start")
	logger.Info("Starting commander")

	actionConfig := new(action.Configuration)
	// Initialize the rest of the actionConfig
	initActionConfig(actionConfig, false)

	//kubeClient := kube.New(kubeConfig)
	//
	//helmClient := helm.NewActionConfig(kubeConfig, nil)
	//prov := kubeProv.New(helmClient, kubeClient)
	//
	//httpServer := api.NewHttp(kubeClient)
	//logger.Info(fmt.Sprintf("Starting HTTP server on port %s", appConfig.HttpPort))
	//httpServer.Serve(appConfig.HttpPort)
	//
	//grpcServer := api.NewGRPC(&prov)
	//logger.Info(fmt.Sprintf("Starting gRPC server on port %s", appConfig.GRPCPort))
	//grpcServer.Serve(appConfig.GRPCPort)
}


func initActionConfig(actionConfig *action.Configuration, allNamespaces bool) {
	kc := kube.New(helm.KubeConfig())

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		// TODO return error
		log.Fatal(err)
	}
	var namespace string
	if !allNamespaces {
		namespace = helm.GetNamespace()
	}

	var store *storage.Storage
	switch os.Getenv("HELM_DRIVER") {
	case "secret", "secrets", "":
		d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
		store = storage.Init(d)
	case "memory":
		d := driver.NewMemory()
		store = storage.Init(d)
	default:
		// Not sure what to do here.
		panic("Unknown driver in HELM_DRIVER: " + os.Getenv("HELM_DRIVER"))
	}

	actionConfig.RESTClientGetter = helm.KubeConfig()
	actionConfig.KubeClient = kc
	actionConfig.Releases = store
}