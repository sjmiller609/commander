package kubernetesProv

import (
	"github.com/sirupsen/logrus"

	"github.com/astronomerio/commander/config"
	"github.com/astronomerio/commander/helm"
	"github.com/astronomerio/commander/kubernetes"
	"github.com/astronomerio/commander/pkg/proto"
	"github.com/astronomerio/commander/utils"

)

var (
	log       = logrus.WithField("package", "kubernetes")
	appConfig = config.Get()
)

// KubeProvisioner is capable of deploying and maintaining jobs on Kubernetes.
type KubeProvisioner struct {
	helm  *helm.Client
	kube  *kubernetes.Client
}

func New(helm *helm.Client, kube *kubernetes.Client) KubeProvisioner {
	return KubeProvisioner{
		helm: helm,
		kube: kube,
	}
}

func (k *KubeProvisioner) InstallDeployment(request *proto.CreateDeploymentRequest) (*proto.CreateDeploymentResponse) {
	response := &proto.CreateDeploymentResponse{
		Deployment: &proto.Deployment{},
	}

	if len(request.Secrets) > 0 {
		for _, secret := range request.Secrets {
			k.kube.Secret.Create(secret.Name, secret.Key, secret.Value, appConfig.KubeNamespace, request.ReleaseName)
		}
	}

	options, err := utils.ParseJSON(request.RawConfig)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response
	}

	install, err := k.helm.InstallRelease(request.ReleaseName, request.Chart.Name, request.Chart.Version, appConfig.KubeNamespace, options)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response
	}
	response.Result = BuildResult(true, "Deployment Created")
	response.Deployment.ReleaseName = install.Release.Name
	return response
}

func (k *KubeProvisioner) UpdateDeployment(request *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse) {
	response := &proto.UpdateDeploymentResponse{
		Deployment: &proto.Deployment{},
	}

	options, err := utils.ParseJSON(request.RawConfig)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response
	}

	update, err := k.helm.UpdateRelease(request.ReleaseName, request.Chart.Name, request.Chart.Version, options)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response
	}

	response.Result = BuildResult(true, "Deployment Updated")
	response.Deployment.ReleaseName = update.Release.Name
	return response
}

func BuildResult(success bool, message string) *proto.Result {
	return &proto.Result {
		Success: success,
		Message: message,
	}
}