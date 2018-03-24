package cmd

import (
	"os"

	"github.com/astronomerio/commander/api"
	"github.com/astronomerio/commander/api/v1"
	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/kubernetes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log       = logrus.WithField("package", "cmd")
	appConfig = config.Get()
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

	// Create new API client and begin accepting requests
	client := api.NewClient()
	initDeploymentRouteHandler(client)
	client.Serve(appConfig.Port)
}

func initDeploymentRouteHandler(client *api.Client) {
	logger := log.WithField("function", "initDeploymentRouteHandler")
	logger.Debug("Entered initDeploymentRouteHandler")

	kubernetesProvisioner, err := kubernetes.NewKubeProvisioner()
	if err != nil {
		logger.Panic(err)
	}

	// Alternate provisioners can be swapped here
	deploymentRouteHandler := v1.NewDeploymentRouteHandler(kubernetesProvisioner)
	client.AppendRouteHandler(deploymentRouteHandler)
}
