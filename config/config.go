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
	Port          string `mapstructure:"PORT"`
	KubeConfig    string `mapstructure:"KUBECONFIG"`
	KubeCoreNS string `mapstructure:"KUBE_CORE_NAMESPACE"`
	KubeAirflowNS string `mapstructure:"KUBE_AIRFLOW_NAMESPACE"`
	KubeClickstreamNS string `mapstructure:"KUBE_CLICKSTREAM_NAMESPACE"`
	HelmRepo	  string `mapstructure:"HELM_REPO"`
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
	viper.SetDefault("PORT", "8880")
	viper.SetDefault("KUBECONFIG", "")
	viper.SetDefault("KUBE_CORE_NAMESPACE", "astronomer")
	viper.SetDefault("KUBE_AIRFLOW_NAMESPACE", "astronomer")
	viper.SetDefault("KUBE_CLICKSTREAM_NAMESPACE", "astronomer")
	viper.SetDefault("HELM_REPO", "https://helm.astronomer.io")
}

// Get returns a populated config struct
func Get() *Configuration {
	return &config
}
