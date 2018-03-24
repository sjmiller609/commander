package helm

import (
	"errors"
	"fmt"
	"path/filepath"
	"os"
	"strings"

	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/repo"
)

func KubeNamespaceForChart(chart string) string {
	switch chart {
	case "airflow":
		return appConfig.KubeAirflowNS
	case "clickstream":
		return appConfig.KubeClickstreamNS
	default:
		return appConfig.KubeCoreNS
	}
}

func (c *Client) AddRepository(cacheFile, repoName, repoUrl string) (*repo.ChartRepository, error) {
	entry := repo.Entry{
		Name:  repoName,
		URL:   repoUrl,
		Cache: cacheFile,
	}
	repository, err := repo.NewChartRepository(&entry, getter.All(c.settings))
	if err != nil {
		return nil, err
	}

	// In this case, the cacheFile is always absolute. So passing empty string is safe.
	err = c.DownloadRepository(repository)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

// Loads the client.repo contents into an index file
func (c *Client) LoadRepoIndex() (error) {
	index, err := repo.LoadIndexFile(c.repo.Config.Cache)
	if err != nil {
		return err
	}
	c.repo.IndexFile = index
	return nil
}

func (c *Client) AcquireChartPath(chart string, version string) (string, error) {
	parts := strings.SplitN(chart, "/", 2)
	if len(parts) != 2 {
		return "", errors.New("chart name should be in the format 'repo/chart name'")
	}
	name := strings.TrimSpace(parts[1])
	version = strings.TrimSpace(version)

	if version == "" {
		latestVersion, err := c.LatestVersion(name)
		if err != nil {
			return "", err
		}
		version = latestVersion
	}

	if !c.ChartDownloaded(name, version) {
		fmt.Println("Chart not downloaded")

		if !c.ChartKnown(name, version) {
			fmt.Println("Chart not known")
			// if not downloaded, we should see if it is in the index.yaml
			// if not download the latest repo
			// after downloading latest repo, check again if its in the index.yaml
			// if not, return error that version is unknown
			// if either it was in the index.yaml, or we've downloaded the latest one and its now there
			//    download the chart
			//    verify the chart exists in the cache now, and then return the path

			c.DownloadRepository(c.repo)
			c.LoadRepoIndex()

			if !c.ChartKnown(name, version) {
				fmt.Println("Chart still not known")
				return "", fmt.Errorf("chart \"%s-%s\" not found in repository", name, version)
			}
		}

		err := c.DownloadChart(chart, version)

		if err != nil {
			return "", err
		}

		if !c.ChartDownloaded(name, version) {
			return "", fmt.Errorf("could not download chart for \"%s-%s\"", name, version)
		}
	}

	return c.buildChartPath(name, version), nil
}

func (c *Client) ChartDownloaded(name string, version string) (bool) {
	absPath := c.buildChartPath(name, version)
	fmt.Println(absPath)
	_, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	return true
}

func (c *Client) ChartKnown(name string, version string) (bool) {
	chart, err := c.repo.IndexFile.Get(name, version)
	if err != nil {
		fmt.Printf("ChartKnown error: %s", err.Error())
		return false
	}
	fmt.Printf("%+v\n", chart)

	return true
}

func (c *Client) LatestVersion(name string) (string, error) {
	chart, err := c.repo.IndexFile.Get(name, "")
	if err != nil {
		fmt.Printf("Unable to find chart by name: %s", err.Error())
		return "", err
	}
	return chart.Version, nil
}

func (c *Client) DownloadRepository(repository *repo.ChartRepository) error {
	// In this case, the cacheFile is always absolute. So passing empty string is safe.
	if err := repository.DownloadIndexFile(""); err != nil {
		return fmt.Errorf("looks like %q is not a valid chart repository or cannot be reached: %s", c.repoUrl, err.Error())
	}
	return nil
}

func (c *Client) DownloadChart(name string, version string) error {
	dl := downloader.ChartDownloader{
		HelmHome: c.settings.Home,
		Out:      os.Stdout,
		Getters:  getter.All(c.settings),
	}
	_, _, err := dl.DownloadTo(name, version, c.settings.Home.Archive())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DefaultNamespace() string {
	if ns, _, err := kube.GetConfig(c.settings.KubeContext).Namespace(); err == nil {
		return ns
	}
	return "default"
}


//func LocateChartPath(repoUrl string, name string, version string) (string, error) {
//	name = strings.TrimSpace(name)
//	version = strings.TrimSpace(version)
//
//	filename, err := GetLocalChart(name)
//	if err != nil {
//		filename, err = DownloadChart(repoUrl, name, version)
//	}
//	return filename, nil
//}
//
//func GetLocalChart(name string) (string, error) {
//	_, err := os.Stat(name)
//	if err != nil {
//		return name, fmt.Errorf("path %q not found", name)
//	}
//
//	abs, err := filepath.Abs(name)
//	if err != nil {
//		return abs, err
//	}
//	return abs, nil
//}
//
//func DownloadChart(repoUrl string, name string, version string) (string, error) {
//	log.Debugf("Downloading chart %s\n", repoUrl)
//	dl := downloader.ChartDownloader{
//		HelmHome: settings.Home,
//		Out:      os.Stdout,
//		Getters:  getter.All(settings),
//	}
//
//	chartUrl, err := repo.FindChartInRepoURL(repoUrl, name, version,"", "", "", getter.All(settings))
//	if err != nil {
//		return "", err
//	}
//	name = chartUrl
//
//	if _, err := os.Stat(settings.Home.Archive()); os.IsNotExist(err) {
//		os.MkdirAll(settings.Home.Archive(), 0744)
//	}
//
//	filename, _, err := dl.DownloadTo(name, version, "./")
//	if err != nil {
//		log.Errorf("Failed to download %s", name)
//		return filename, err
//	}
//	return chartUrl, nil
//
//	//lname, err := filepath.Abs(filename)
//	//if err != nil {
//	//	return filename, err
//	//}
//	//log.Debugf("Fetched %s to %s\n", name, filename)
//	//return lname, nil
//}

//
//func (c *Client) ChartVersionValid(chart string, version string) bool {
//	insp.output = chartOnly
//
//	cp, err := chartutil.locateChartPath(c.repo, chart, version, insp.verify, insp.keyring,
//		insp.certFile, insp.keyFile, insp.caFile)
//	if err != nil {
//		return err
//	}
//	insp.chartpath = cp
//
//	chrt, err := chartutil.Load(i.chartpath)
//	if err != nil {
//		return err
//	}
//	cf, err := yaml.Marshal(chrt.Metadata)
//	if err != nil {
//		return err
//	}
//
//	if i.output == chartOnly || i.output == both {
//		fmt.Fprintln(i.out, string(cf))
//	}
//
//	if (i.output == valuesOnly || i.output == both) && chrt.Values != nil {
//		if i.output == both {
//			fmt.Fprintln(i.out, "---")
//		}
//		fmt.Fprintln(i.out, chrt.Values.Raw)
//	}
//
//	return nil
//
//	return false
//}

// ensureDirectories checks to see if $HELM_HOME exists.
// If $HELM_HOME does not exist, this function will create it and all its subdirectories
func (c *Client) ensureDirectories() error {
	configDirectories := []string{
		c.settings.Home.String(),
		c.settings.Home.Repository(),
		c.settings.Home.Cache(),
		c.settings.Home.LocalRepository(),
		c.settings.Home.Plugins(),
		c.settings.Home.Starters(),
		c.settings.Home.Archive(),
	}
	fmt.Println(configDirectories)
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			if err := os.MkdirAll(p, 0755); err != nil {
				return err
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("%s must be a directory", p)
		}
	}

	return nil
}

// ensures that the astronomer-ee chart repo is registered with helm
// if not it will get added
func (c *Client) ensureAstroRepo() error {
	repoFile := c.settings.Home.RepositoryFile()
	if fi, err := os.Stat(repoFile); err != nil {
		f := repo.NewRepoFile()
		stableRepo, err := c.AddRepository(c.localRepoPath(stableRepository), stableRepository, stableRepositoryURL)
		if err != nil {
			return err
		}
		astroRepo, err := c.AddRepository(c.localRepoPath(astroRepoName), astroRepoName, c.repoUrl)
		if err != nil {
			return err
		}
		f.Add(stableRepo.Config)
		f.Add(astroRepo.Config)
		c.repo = astroRepo
		c.LoadRepository(astroRepo.Config.Cache, "", "")

		if err := f.WriteFile(repoFile, 0644); err != nil {
			return err
		}

		err = c.LoadRepoIndex()
		if err != nil {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", repoFile)
	}
	return nil
}

func (c *Client) buildChartPath(name string, version string) string {
	filename := fmt.Sprintf("%s-%s.tgz", name, version)
	absPath := filepath.Join(c.settings.Home.Archive(), filename)
	return absPath
}

func (c *Client) localRepoPath(repoName string) string {
	return c.settings.Home.CacheIndex(repoName)
}

// Seems that the repo.Load() function assumes the repo Entry name is a directory, which doesn't
// seem to be the case anymore (or yet)

func (c *Client) LoadRepository(cacheFile, repoName, repoUrl string) (*repo.ChartRepository, error) {
	entry := repo.Entry{
		Name:  repoName,
		URL:   repoUrl,
		Cache: cacheFile,
	}
	repository, err := repo.NewChartRepository(&entry, getter.All(c.settings))
	if err != nil {
		return nil, err
	}

	err = repository.Load()
	if err != nil {
		return nil, err
	}
	return repository, nil
}