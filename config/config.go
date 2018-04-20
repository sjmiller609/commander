package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config Configuration
	prefix = "COMMANDER"
	log    = logrus.WithField("package", "config")
)

// Configuration is a struct to hold provisioner configs
type Configuration struct {
	DebugMode     bool   `mapstructure:"DEBUG_MODE"`
	HttpPort          string `mapstructure:"HTTP_PORT"`
	GRPCPort          string `mapstructure:"GRPC_PORT"`
	KubeConfig    string `mapstructure:"KUBECONFIG"`
	KubeNamespace string `mapstructure:"KUBE_NAMESPACE"`
	HelmRepo	  string `mapstructure:"HELM_REPO"`
	HelmRepoName	  string `mapstructure:"HELM_REPO_NAME"`
	TillerHost		string `mapstructure:"TILLER_HOST"`
}

// Log will log the configuation struct out
func (c *Configuration) Log() {
	logger := log.WithField("function", "init")
	logger.Info(fmt.Sprintf("%+v", c))
}

// Initalize configuration
func Init() {
	logger := log.WithField("function", "init")
	logger.Debug("Initializing configuration")

	setDefaults()
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		logger.Panic(fmt.Sprintf("Unable to decode configuration, %v", err))
	}
}

// Set some default values
func setDefaults() {
	viper.SetDefault("DEBUG_MODE", true)
	viper.SetDefault("HTTP_PORT", "8880")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("KUBECONFIG", "")
	viper.SetDefault("KUBE_NAMESPACE", "astronomer")
	viper.SetDefault("HELM_REPO", "https://helm.astronomer.io")
	viper.SetDefault("HELM_REPO_NAME", "astronomer-ee")
	viper.SetDefault("TILLER_HOST", "127.0.0.1:34477")
}

// Get returns a populated config struct
func Get() *Configuration {
	return &config
}
