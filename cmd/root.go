package cmd

import (
	"fmt"
	"os"

	"github.com/astronomerio/commander/pkg/proto"

	"github.com/astronomerio/commander/api"
	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/kubernetes"
	kubeProv "github.com/astronomerio/commander/provisioner/kubernetes"
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

	kubeConfig, err := kubernetes.GetKubeConfig()
	if err != nil {
		logger.Panic(err)
	}

	kubeClient, err := kubernetes.New(kubeConfig)
	if err != nil {
		logger.Panic(err)
	}

	//err = kubeClient.Namespace.Ensure(appConfig.KubeNamespace)
	//if err != nil {
	//	logger.Panic(err)
	//}

	helmClient := helm.New(kubeClient, appConfig.HelmRepo)

	prov := kubeProv.New(helmClient, kubeClient)

	// TODO: remove when done testing
	request := &proto.DeleteDeploymentRequest{
		ReleaseName: "dangling-narwhal",
	}
	response, err := prov.DeleteDeployment(request)
	fmt.Println("helloooooooooo")
	fmt.Println(fmt.Sprintf("%v", response))
	fmt.Println(fmt.Sprintf("%v", err))

	httpServer := api.NewHttp()
	logger.Info(fmt.Sprintf("Starting HTTP server on port %s", appConfig.HttpPort))
	httpServer.Serve(appConfig.HttpPort)

	grpcServer := api.NewGRPC(&prov)
	logger.Info(fmt.Sprintf("Starting gRPC server on port %s", appConfig.GRPCPort))
	grpcServer.Serve(appConfig.GRPCPort)
}
