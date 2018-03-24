package cmd

import (
	"fmt"
	"os"

	"github.com/astronomerio/commander/api"
	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/kubernetes"
	"github.com/astronomerio/commander/helm"
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

	helmClient := helm.New(appConfig.HelmRepo)

	fmt.Println("DERP Helm initialized lets install a release")
	_, err := helmClient.InstallRelease("astronomer-ee/airflow", "",      appConfig.KubeAirflowNS, map[string]interface{}{})
	fmt.Println(err)
	//_, err = helmClient.InstallRelease("astronomer-ee/airflow", "0.1.2", appConfig.KubeAirflowNS, map[string]interface{}{})
	//fmt.Println(err)

	kubeConfig, err := kubernetes.GetKubeConfig()
	if err != nil {
		logger.Panic(err)
	}

	kubeClient, err := kubernetes.New(kubeConfig)
	if err != nil {
		logger.Panic(err)
	}

	// Ensure all namespaces exist
	kubeClient.Namespace.Ensure(appConfig.KubeCoreNS)
	kubeClient.Namespace.Ensure(appConfig.KubeAirflowNS)
	kubeClient.Namespace.Ensure(appConfig.KubeClickstreamNS)

	// Create new API client and begin accepting requests
	server := api.NewServer()
	logger.Info(fmt.Sprintf("Starting gRPC server on port %s", appConfig.Port))
	server.Serve(appConfig.Port)
}

func getProvisioner() {
	//kubernetesProvisioner, err := kubernetes.NewKubeProvisioner()
	//if err != nil {
	//	logger.Panic(err)
	//}
	//
	//// Alternate provisioners can be swapped here
	//client.AppendRouteHandler(deploymentRouteHandler)
}