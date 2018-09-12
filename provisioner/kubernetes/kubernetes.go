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
	helm *helm.Client
	kube *kubernetes.Client
}

func New(helm *helm.Client, kube *kubernetes.Client) KubeProvisioner {
	return KubeProvisioner{
		helm: helm,
		kube: kube,
	}
}

func (k *KubeProvisioner) InstallDeployment(request *proto.CreateDeploymentRequest) (*proto.CreateDeploymentResponse, error) {
	response := &proto.CreateDeploymentResponse{
		Deployment: &proto.Deployment{},
	}

	if len(request.Secrets) > 0 {
		for _, secret := range request.Secrets {
			k.kube.Secret.Create(secret.Name, secret.Data, request.Namespace, request.ReleaseName)
		}
	}

	options, err := utils.ParseJSON(request.RawConfig)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	install, err := k.helm.InstallRelease(request.ReleaseName, request.Chart.Name, request.Chart.Version, request.Namespace, options)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}
	response.Result = BuildResult(true, "Deployment Created")
	response.Deployment.ReleaseName = install.Release.Name
	return response, nil
}

func (k *KubeProvisioner) UpdateDeployment(request *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse, error) {
	response := &proto.UpdateDeploymentResponse{
		Deployment: &proto.Deployment{},
	}

	options, err := utils.ParseJSON(request.RawConfig)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	update, err := k.helm.UpdateRelease(request.ReleaseName, request.Chart.Name, request.Chart.Version, options)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	response.Result = BuildResult(true, "Deployment Updated")
	response.Deployment.ReleaseName = update.Release.Name
	return response, nil
}

func (k *KubeProvisioner) DeleteDeployment(request *proto.DeleteDeploymentRequest) (*proto.DeleteDeploymentResponse, error) {
	response := &proto.DeleteDeploymentResponse{
		Deployment: &proto.Deployment{},
	}

	// (postgres delete will happen on houston)

	releaseName, info, err := k.helm.DeleteRelease(request.ReleaseName)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	// secret delete (test)
	err = k.kube.Secret.DeleteByRelease(request.ReleaseName, request.Namespace)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	// PVCs delete
	err = k.kube.PersistentVolumeClaim.DeleteByRelease(request.ReleaseName, request.Namespace)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	if request.DeleteNamespace {
		// Namespace delete
		err = k.kube.Namespace.Delete(request.Namespace)
		if err != nil {
			response.Result = BuildResult(false, err.Error())
			return response, nil
		}
	}

	response.Result = BuildResult(true, "Deployment Deleted")
	response.Deployment.ReleaseName = releaseName
	response.Deployment.Info = info

	return response, nil
}

func (k *KubeProvisioner) GetSecret(request *proto.GetSecretRequest) (*proto.GetSecretResponse, error) {
	response := &proto.GetSecretResponse{
		Secret: &proto.Secret{},
	}
	secret, err := k.kube.Secret.Get(request.Name, request.Namespace)
	if err != nil {
		response.Result = BuildResult(false, err.Error())
		return response, nil
	}

	response.Result = BuildResult(true, "")
	response.Secret.Name = request.Name
	response.Secret.Data = secret.StringData
	return response, nil
}

func BuildResult(success bool, message string) *proto.Result {
	return &proto.Result{
		Success: success,
		Message: message,
	}
}
