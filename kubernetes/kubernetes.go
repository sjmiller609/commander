package kubernetes

import (
	"fmt"

	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/provisioner"
	"github.com/sirupsen/logrus"
	apiappsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	// Required to authenticate against GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	log       = logrus.WithField("package", "kubernetes")
	appConfig = config.Get()
)

// KubeProvisioner is capable of deploying and maintaining jobs on Kubernetes.
type KubeProvisioner struct {
	deploymentsClient typedappsv1beta2.DeploymentInterface
	stsClient         typedappsv1beta2.StatefulSetInterface
}

// NewKubeProvisioner returns a new KubeProvisioner
func NewKubeProvisioner() (*KubeProvisioner, error) {
	logger := log.WithField("function", "NewKubeProvisioner")
	logger.Debug("Creating Kubernetes provisioner")

	var (
		kubeConfig *rest.Config
		configErr  error
	)

	if appConfig.KubeConfig != "" {
		logger.Debug("Using config at ", appConfig.KubeConfig)
		kubeConfig, configErr = clientcmd.BuildConfigFromFlags("", appConfig.KubeConfig)
	} else {
		logger.Debug("Using in-cluster config")
		kubeConfig, configErr = rest.InClusterConfig()
	}

	if configErr != nil {
		return nil, configErr
	}

	clientset, clientErr := kubernetes.NewForConfig(kubeConfig)
	if clientErr != nil {
		return nil, clientErr
	}

	return &KubeProvisioner{
		deploymentsClient: clientset.AppsV1beta2().Deployments(appConfig.KubeNamespace),
		stsClient:         clientset.AppsV1beta2().StatefulSets(appConfig.KubeNamespace),
	}, nil
}

// ListDeployments returns a list of known deployments.
func (p *KubeProvisioner) ListDeployments(organizationID string) (*provisioner.ListDeploymentResponse, error) {
	logger := log.WithField("function", "ListDeployments")
	logger.Debug("Entered ListDeployments")

	labels := fmt.Sprintf("organization=%s, tier=airflow-core", organizationID)
	deployments, listErr := p.deploymentsClient.List(metav1.ListOptions{
		LabelSelector: labels,
	})

	if listErr != nil {
		return nil, listErr
	}

	// Use a map to create a unique list
	items := make(map[string]bool)
	for _, deployment := range deployments.Items {
		if release, ok := deployment.ObjectMeta.Labels["release"]; ok {
			items[release] = true
		}
	}

	releaseNames := []string{}
	for key := range items {
		releaseNames = append(releaseNames, key)
	}

	resp := &provisioner.ListDeploymentResponse{
		Items: releaseNames,
	}
	return resp, nil
}

// PatchDeployment patches a deployment in place.
func (p *KubeProvisioner) PatchDeployment(deploymentID string, req *provisioner.PatchDeploymentRequest) (*provisioner.PatchDeploymentResponse, error) {
	logger := log.WithField("function", "PatchDeployment")
	logger.Debug("Entered PatchDeployment")

	webserverName := fmt.Sprintf("%s-webserver", deploymentID)
	_, webErr := p.patchDeployment(webserverName, req.Image)
	if webErr != nil {
		return nil, webErr
	}

	schedulerName := fmt.Sprintf("%s-scheduler", deploymentID)
	_, schedErr := p.patchDeployment(schedulerName, req.Image)
	if schedErr != nil {
		return nil, schedErr
	}

	flowerName := fmt.Sprintf("%s-flower", deploymentID)
	_, flowerErr := p.patchDeployment(flowerName, req.Image)
	if flowerErr != nil {
		return nil, flowerErr
	}

	workerName := fmt.Sprintf("%s-worker", deploymentID)
	_, workerErr := p.patchStatefulSet(workerName, req.Image)
	if workerErr != nil {
		return nil, workerErr
	}

	resp := &provisioner.PatchDeploymentResponse{}
	return resp, nil
}

// Patches a Kubernetes Deployment resource.
func (p *KubeProvisioner) patchDeployment(name, image string) (*apiappsv1beta2.Deployment, error) {
	logger := log.WithField("function", "patchDeployment")
	logger.Debug("Entered patchDeployment")

	// Get the deployment using the given name.
	d, getErr := p.deploymentsClient.Get(name, metav1.GetOptions{})
	if getErr != nil {
		return d, getErr
	}

	// Update the deployment.
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		d.Spec.Template.Spec.Containers[0].Image = image
		_, updateErr := p.deploymentsClient.Update(d)
		return updateErr
	})

	return d, retryErr
}

// Patches a Kubernetes StatefulSet resource.
func (p *KubeProvisioner) patchStatefulSet(name, image string) (*apiappsv1beta2.StatefulSet, error) {
	logger := log.WithField("function", "patchStatefulSet")
	logger.Debug("Entered patchStatefulSet")

	// Get the statefulset using the given name.
	s, getErr := p.stsClient.Get(name, metav1.GetOptions{})
	if getErr != nil {
		return s, getErr
	}

	// Update the statefulset.
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		s.Spec.Template.Spec.Containers[0].Image = image
		s.Spec.Template.Spec.Containers[1].Image = image
		_, updateErr := p.stsClient.Update(s)
		return updateErr
	})

	return s, retryErr
}
