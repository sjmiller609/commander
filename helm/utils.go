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
	logger := log.WithField("function", "AcquireChartPath")
	logger.Debugf("AcquireChartPath(%s, %s)", chart, version)

	parts := strings.SplitN(chart, "/", 2)
	if len(parts) != 2 {
		logger.Debug("Invalid chart name")
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
		logger.Debug("Chart not downloaded")

		if !c.ChartKnown(name, version) {
			logger.Debug("Chart not known by repository")
			// if not downloaded, we should see if it is in the index.yaml
			// if not download the latest repo
			// after downloading latest repo, check again if its in the index.yaml
			// if not, return error that version is unknown
			// if either it was in the index.yaml, or we've downloaded the latest one and its now there
			//    download the chart
			//    verify the chart exists in the cache now, and then return the path

			logger.Debug("Update repo")
			c.DownloadRepository(c.repo)
			c.LoadRepoIndex()

			if !c.ChartKnown(name, version) {
				logger.Debug("Chart still unknown")
				return "", fmt.Errorf("chart \"%s-%s\" not found in repository", name, version)
			}
		}

		err := c.DownloadChart(chart, version)

		if err != nil {
			return "", err
		}

		if !c.ChartDownloaded(name, version) {
			logger.Error("Error downloading chart")

			return "", fmt.Errorf("could not download chart for \"%s-%s\"", name, version)
		}
	}

	return c.buildChartPath(name, version), nil
}

func (c *Client) ChartDownloaded(name string, version string) (bool) {
	absPath := c.buildChartPath(name, version)
	_, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	return true
}

func (c *Client) ChartName(chartName string) string {
	return fmt.Sprintf("%s/%s", appConfig.HelmRepoName, chartName)
}

func (c *Client) ChartKnown(name string, version string) (bool) {
	_, err := c.repo.IndexFile.Get(name, version)
	if err != nil {
		fmt.Printf("ChartKnown error: %s", err.Error())
		return false
	}
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
	astroRepoName := appConfig.HelmRepoName
	repoFile := c.settings.Home.RepositoryFile()
	if fi, err := os.Stat(repoFile); err != nil {
		f := repo.NewRepoFile()
		stableRepo, err := c.AddRepository(c.localRepoPath(stableRepository), stableRepository, stableRepositoryURL)
		if err != nil {
			return err
		}
		c.DownloadRepository(stableRepo)
		if err != nil {
			return err
		}

		astroRepo, err := c.AddRepository(c.localRepoPath(astroRepoName), astroRepoName, c.repoUrl)
		if err != nil {
			return err
		}
		c.DownloadRepository(astroRepo)
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
	} else {
		astroRepo, err := c.AddRepository(c.localRepoPath(astroRepoName), astroRepoName, c.repoUrl)
		if err != nil {
			return err
		}
		c.repo = astroRepo
		c.LoadRepoIndex()
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