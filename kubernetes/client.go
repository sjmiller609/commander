package kubernetes

import (
	"github.com/sirupsen/logrus"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/astronomerio/commander/config"
)

var (
	log       = logrus.WithField("package", "kubernetes")
	appConfig = config.Get()
)

// KubeProvisioner is capable of deploying and maintaining jobs on Kubernetes.
type Client struct {
	clientSet *kube.Clientset
	Namespace *Namespace
}

// NewKubeProvisioner returns a new KubeProvisioner
func New(kubeConfig *rest.Config) (*Client, error) {
	logger := log.WithField("function", "NewKubeProvisioner")
	logger.Debug("Creating Kubernetes client")

	clientSet, clientErr := kube.NewForConfig(kubeConfig)
	if clientErr != nil {
		return nil, clientErr
	}

	return &Client{
		clientSet: clientSet,
		Namespace: &Namespace{
			clientSet: clientSet,
		},
	}, nil
}