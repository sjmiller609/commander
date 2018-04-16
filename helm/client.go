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
	astroRepoName = "astronomer-ee"
	log       = logrus.WithField("package", "kubernetes")
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
		helm.ConnectTimeout(3),
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
func (c *Client) InstallRelease(chartName, chartVersion, namespace string, options map[string]interface{}) (*services.InstallReleaseResponse, error) {
	// the helm pkg client was designed to go out of scope every command, since we don't do that, we need to reset it
	defer c.Reset()

	optionsYaml, err := yaml.Marshal(options)
	if err != nil {
		return nil, err
	}

	chartPath, err := c.AcquireChartPath(chartName, chartVersion)
	if err != nil {
		return nil, err
	}

	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return c.helm.InstallReleaseFromChart(chart,
		namespace,
		helm.ValueOverrides(optionsYaml),
		helm.ReleaseName(""),
		helm.InstallDryRun(false),
		helm.InstallReuseName(false),
		helm.InstallDisableHooks(false),
		helm.InstallTimeout(300),
		helm.InstallWait(false),
	)
}

// update settings of an existing release
func (c *Client) UpdateRelease(releaseName, options map[string]interface{}) (*services.UpdateReleaseResponse, error) {
	return nil, nil
}

// upgrade a release to a later version of the chart
func (c *Client) UpgradeRelease(releaseName, chartVersion string) {

}

// delete a release
func (c *Client) DeleteRelease(releaseName string) {
	c.helm.DeleteRelease(releaseName)
}

// get release status
func (c *Client) FetchRelease() {

}

