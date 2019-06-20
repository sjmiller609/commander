package helm

import (
	"github.com/astronomer/commander/config"
	"github.com/astronomer/commander/kubernetes"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/chartutil"
	"helm.sh/helm/pkg/cli"
	"helm.sh/helm/pkg/kube"
	"helm.sh/helm/pkg/release"
	"helm.sh/helm/pkg/repo"
)

var (
	appConfig = config.Get()
	log       = logrus.WithField("package", "helm")
	stableRepository         = "stable"
	stableRepositoryURL = "https://kubernetes-charts.storage.googleapis.com"
)

type Client struct {
	helm *kube.Client
	repo *repo.ChartRepository
	repoUrl string
	settings cli.EnvSettings
	kubeClient *kubernetes.Client
}

func New(kubeClient *kubernetes.Client, repo string) *Client {
	// create settings object
	flags := pflag.NewFlagSet("production", pflag.PanicOnError)
	settings := cli.EnvSettings{}
	settings.AddFlags(flags)
	settings.Init(flags)

	client := &Client{
		repoUrl: repo,
		settings: settings,
		kubeClient: kubeClient,
	}

	// create helm client
	client.helm = kube.New(kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace))

	// some helm commands expect `helm init` to have happened.
	// as this isn't an exposed function, we'll just manually do the setup we need

	// As of now, the parts of helm init we need are doing are
	// - Creating helm home and all its subdirectories
	// - Preloading repositories for charts
	if err := client.ensureDirectories(); err != nil {
		panic(err.Error())
	}
	if err := client.ensureAstroRepo(); err != nil {
		panic(err.Error())
	}

	return client
}

func (c *Client) Reset() {
	c.helm = kube.New(nil)
}

// install a new chart release
func (c *Client) InstallRelease(releaseName, chartName, chartVersion, namespace string, options map[string]interface{}) (*release.Release, error) {
	logger := log.WithField("function", "InstallRelease")

	// the helm pkg client was designed to go out of scope every command, since we don't do that, we need to reset it
	defer c.Reset()

	optionsYaml, err := yaml.Marshal(options)
	if err != nil {
		return nil, err
	}

	chartPath, err := c.AcquireChartPath(c.ChartName(chartName), chartVersion)
	if err != nil {
		logger.Errorf("#AcquireChartPath: %s", err.Error())
		return nil, err
	}

	chartPath, err := chartutil.LoadChartfile(chartPath)
	if err != nil {
		return nil, err
	}

	logger.Debug("helm#InstallReleaseFromChart")
	client := action.NewInstall(cfg)
	client.Namespace = namespace
	client.DryRun = false
	client.DisableHooks = false
	client.Timeout = 300
	client.Wait = false
	//client.ValueOptions
	return client.Run(chartPath)
}

// update settings of an existing release
func (c *Client) UpdateRelease(releaseName, chartName, chartVersion string, options map[string]interface{}) (*services.UpdateReleaseResponse, error) {
	logger := log.WithField("function", "UpdateRelease")

	optionsYaml, err := yaml.Marshal(options)
	if err != nil {
		return nil, err
	}

	chartPath, err := c.AcquireChartPath(c.ChartName(chartName), chartVersion)
	if err != nil {
		logger.Errorf("#AcquireChartPath: %s", err.Error())
		return nil, err
	}

	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	logger.Debug("helm#UpdateReleaseFromChart")
	return c.helm.UpdateReleaseFromChart(releaseName,
		chart,
		helm.UpdateValueOverrides(optionsYaml),
		helm.UpgradeDryRun(false),
		helm.ReuseValues(true),
		helm.UpgradeDisableHooks(false),
		helm.UpgradeTimeout(300),
		helm.UpgradeWait(false),
	)

	return nil, nil
}

// upgrade a release to a later version of the chart
func (c *Client) UpgradeRelease(releaseName, chartName, chartVersion string, options map[string]interface{}) (*services.UpdateReleaseResponse, error) {
	logger := log.WithField("function", "UpgradeRelease")

	optionsYaml, err := yaml.Marshal(options)
	if err != nil {
		return nil, err
	}

	chartPath, err := c.AcquireChartPath(c.ChartName(chartName), chartVersion)
	if err != nil {
		logger.Errorf("#AcquireChartPath: %s", err.Error())
		return nil, err
	}

	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	logger.Debug("helm#UpgradeReleaseFromChart")
	return c.helm.UpdateReleaseFromChart(releaseName,
		chart,
		helm.UpdateValueOverrides(optionsYaml),
		helm.UpgradeDryRun(false),
		helm.ReuseValues(true),
		helm.UpgradeTimeout(300),
		helm.UpgradeWait(false),
	)

	return nil, nil
}

// delete a release
func (c *Client) DeleteRelease(releaseName string) (string, string, error) {
	logger := log.WithField("function", "DeleteRelease")

	logger.Debug("helm#DeleteRelease")
	response, err := c.helm.DeleteRelease(releaseName, helm.DeletePurge(true))
	if err != nil {
		return "", "", err
	}
	return response.GetRelease().GetName(), response.GetInfo(), nil
}

// get release status
func (c *Client) FetchRelease() {

}

