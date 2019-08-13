package helm

import (
	"github.com/sirupsen/logrus"
	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/chart/loader"
	"helm.sh/helm/pkg/cli"
	"helm.sh/helm/pkg/downloader"
	"helm.sh/helm/pkg/getter"
	"helm.sh/helm/pkg/release"
)


var (
	log       = logrus.WithField("package", "helm")
	settings   cli.EnvSettings
)

type Client struct {
	actionConfig *action.Configuration
}


func NewClient(actionConfig *action.Configuration) *Client {
	client := &Client{
		actionConfig: actionConfig,
	}
	return client
}


// install a new chart release
func (c *Client) InstallRelease(releaseName, chart, chartVersion, namespace string, options map[string]interface{}) (*release.Release, error) {
	client := action.NewInstall(c.actionConfig)
	logger := log.WithField("function", "InstallRelease")

	logger.Info("Starting commander")
	logger.Info("Original chart version: %q", client.Version)

	if client.Version == "" && client.Devel {
		logger.Info("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	client.ReleaseName = releaseName

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	logger.Info("CHART PATH: %s\n", cp)

	vals := make(map[string]interface{})

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					//Out:        out,
					ChartPath:  cp,
					Keyring:    client.ChartPathOptions.Keyring,
					SkipUpdate: false,
					Getters:    getter.All(settings),
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = GetNamespace()
	return client.Run(chartRequested, vals)
}

func (c *Client) UpdateRelease(s string, s2 string, s3 string, options map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (c *Client) UpgradeRelease(s string, s2 string, s3 string, options map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (c *Client) DeleteRelease(s string) (interface{}, interface{}, interface{}) {
	return nil, nil, nil
}
