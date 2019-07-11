package cmd

import (
	"github.com/astronomer/commander/kubernetes"
	//"helm.sh/helm/pkg/kube"
	"os"

	//"github.com/astronomer/commander/api"
	"github.com/astronomer/commander/config"
	//"github.com/astronomer/commander/helm"
	//"github.com/astronomer/commander/kubernetes"
	//kubeProv "github.com/astronomer/commander/provisioner/kubernetes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log       = logrus.WithField("package", "cmd")
	appConfig = config.Get()

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

	config.Init()

	if appConfig.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logger := log.WithField("function", "start")
	logger.Info("Starting commander")

	kubeConfig := kubernetes.GetKubeConfig()
	_ = kubeConfig

	//kubeClient := kube.New(kubeConfig)

	//helmClient := helm.NewActionConfig(kubeConfig, kubeClient)
	//prov := kubeProv.New(helmClient, kubeClient)

	//httpServer := api.NewHttp(kubeClient)
	//logger.Info(fmt.Sprintf("Starting HTTP server on port %s", appConfig.HttpPort))
	//httpServer.Serve(appConfig.HttpPort)
	//
	//grpcServer := api.NewGRPC(&prov)
	//logger.Info(fmt.Sprintf("Starting gRPC server on port %s", appConfig.GRPCPort))
	//grpcServer.Serve(appConfig.GRPCPort)
}
