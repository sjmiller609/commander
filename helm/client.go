package helm

import (
	"github.com/sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/spf13/pflag"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/proto/hapi/services"

	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/kubernetes"
)

var (
	appConfig = config.Get()
	log       = logrus.WithField("package", "helm")
	stableRepository         = "stable"
	stableRepositoryURL = "https://kubernetes-charts.storage.googleapis.com"
)

type Client struct {
	helm *helm.Client
	helmOptions []helm.Option
	repo *repo.ChartRepository
	repoUrl string
	settings environment.EnvSettings
	kubeClient *kubernetes.Client
	tillerTunnel *kube.Tunnel
}

func New(kubeClient *kubernetes.Client, repo string) *Client {
	// create settings object
	flags := pflag.NewFlagSet("production", pflag.PanicOnError)
	settings := environment.EnvSettings{}
	settings.AddFlags(flags)
	settings.Init(flags)

	client := &Client{
		repoUrl: repo,
		settings: settings,
		kubeClient: kubeClient,
	}

	// open tunnel to tiller (if needed)
	client.OpenTunnel()

	client.helmOptions = []helm.Option{
		helm.ConnectTimeout(5),
		helm.Host(client.settings.TillerHost),
	}

	// create helm client
	client.helm = helm.NewClient(client.helmOptions...)

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
	c.helm = helm.NewClient(c.helmOptions...)
}

// install a new chart release
func (c *Client) InstallRelease(releaseName, chartName, chartVersion, namespace string, options map[string]interface{}) (*services.InstallReleaseResponse, error) {
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

	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	logger.Debug("helm#InstallReleaseFromChart")
	return c.helm.InstallReleaseFromChart(chart,
		namespace,
		helm.ValueOverrides(optionsYaml),
		helm.ReleaseName(releaseName),
		helm.InstallDryRun(false),
		helm.InstallReuseName(false),
		helm.InstallDisableHooks(false),
		helm.InstallTimeout(300),
		helm.InstallWait(false),
	)
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
		helm.UpgradeRecreate(true),
		helm.UpgradeDisableHooks(false),
		helm.UpgradeTimeout(300),
		helm.UpgradeWait(false),
	)

	return nil, nil
}

// upgrade a release to a later version of the chart
func (c *Client) UpgradeRelease(releaseName, chartVersion string) {

}

// delete a release
func (c *Client) DeleteRelease(releaseName string) (string, string, error) {
	logger := log.WithField("function", "DeleteRelease")

	logger.Debug("helm#DeleteRelease")
	response, err := c.helm.DeleteRelease(releaseName)
	if err != nil {
		return "", "", err
	}
	return response.GetRelease().GetName(), response.GetInfo(), nil
}

// get release status
func (c *Client) FetchRelease() {

}

