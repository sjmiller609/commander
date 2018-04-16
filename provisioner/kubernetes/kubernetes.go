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
	response := &proto.CreateDeploymentResponse{}

	options, err := utils.ParseJSON(request.RawConfig)
	if err != nil {
		response.Result = BuildResult(err)
		return response
	}

	install, err := k.helm.InstallRelease(request.Chart.Name, request.Chart.Version, appConfig.KubeNamespace, options)
	if err != nil {
		response.Result = BuildResult(err)
		return response
	}

	_ = install

	return response
}

func (k *KubeProvisioner) UpdateDeployment(request *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse) {
	response := &proto.UpdateDeploymentResponse{}

	// options, err := utils.ParseJSON(request.RawConfig)
	// if err != nil {
	// 	response.Result = BuildResult(err)
	// 	return response
	// }
	//
	// update, err := k.helm.UpdateRelease(request.ReleaseName, options)
	// if err != nil {
	// 	response.Result = BuildResult(err)
	// 	return response
	// }

	// _ = update

	return response
}

// func (k *KubeProvisioner) UpdateDeployment() {
//
// }
// func (k *KubeProvisioner) UpgradeDeployment() {
//
// }
// func (k *KubeProvisioner) DeleteDeployment() {
//
// }
// func (k *KubeProvisioner) FetchDeployments() {
//
// }
//
// func (k *KubeProvisioner) PatchDeployment(patchReq *provisioner.PatchDeploymentRequest) (*provisioner.PatchDeploymentResponse, error) {
// 	//	logger := log.WithField("function", "PatchDeployment")
// 	//	logger.Debug("Entered PatchDeployment")
// 	//	logger.Debug(fmt.Sprintf("%+v", patchReq))
// 	//
// 	//	metadata := patchReq.Metadata
// 	//	labels := fmt.Sprintf("release=%s, tier=%s", metadata.DeploymentID, metadata.ComponentID)
// 	//	depErr := p.patchDeployment(labels, patchReq.Image)
// 	//	if depErr != nil {
// 	//		return nil, depErr
// 	//	}
// 	//
// 	//	stsErr := p.patchStatefulSet(labels, patchReq.Image)
// 	//	if stsErr != nil {
// 	//		return nil, stsErr
// 	//	}
// 	//
// 	//	return &provisioner.PatchDeploymentResponse{}, nil
//
// 	return &provisioner.PatchDeploymentResponse{}, nil
// }
// //// NewKubeProvisioner returns a new KubeProvisioner
// //func NewKubeProvisioner() (*KubeProvisioner, error) {
// //	logger := log.WithField("function", "NewKubeProvisioner")
// //	logger.Debug("Creating Kubernetes provisioner")
// //
// //	var (
// //		kubeConfig *rest.Config
// //		configErr  error
// //	)
// //
// //	if appConfig.KubeConfig != "" {
// //		logger.Debug("Using config at ", appConfig.KubeConfig)
// //		kubeConfig, configErr = clientcmd.BuildConfigFromFlags("", appConfig.KubeConfig)
// //	} else {
// //		logger.Debug("Using in-cluster config")
// //		kubeConfig, configErr = rest.InClusterConfig()
// //	}
// //
// //	if configErr != nil {
// //		return nil, configErr
// //	}
// //
// //	clientset, clientErr := kubernetes.NewForConfig(kubeConfig)
// //	if clientErr != nil {
// //		return nil, clientErr
// //	}
// //
// //	return &KubeProvisioner{
// //		deploymentsClient: clientset.AppsV1beta2().Deployments(appConfig.KubeNamespace),
// //		stsClient:         clientset.AppsV1beta2().StatefulSets(appConfig.KubeNamespace),
// //	}, nil
// //}
//
// // ListDeployments returns a list of known deployments.
// // func (p *KubeProvisioner) ListDeployments(organizationID string) (*provisioner.ListDeploymentResponse, error) {
// // 	logger := log.WithField("function", "ListDeployments")
// // 	logger.Debug("Entered ListDeployments")
//
// // 	labels := fmt.Sprintf("organization=%s, tier=airflow-core", organizationID)
// // 	deployments, listErr := p.deploymentsClient.List(metav1.ListOptions{
// // 		LabelSelector: labels,
// // 	})
//
// // 	if listErr != nil {
// // 		return nil, listErr
// // 	}
//
// // 	// Use a map to create a unique list
// // 	items := make(map[string]bool)
// // 	for _, deployment := range deployments.Items {
// // 		if release, ok := deployment.ObjectMeta.Labels["release"]; ok {
// // 			items[release] = true
// // 		}
// // 	}
//
// // 	releaseNames := []string{}
// // 	for key := range items {
// // 		releaseNames = append(releaseNames, key)
// // 	}
//
// // 	resp := &provisioner.ListDeploymentResponse{
// // 		Items: releaseNames,
// // 	}
// // 	return resp, nil
// // }
// //
// //// PatchDeployment patches a deployment in place.
// //func (p *KubeProvisioner) PatchDeployment(patchReq *provisioner.PatchDeploymentRequest) (*provisioner.PatchDeploymentResponse, error) {
// //	logger := log.WithField("function", "PatchDeployment")
// //	logger.Debug("Entered PatchDeployment")
// //	logger.Debug(fmt.Sprintf("%+v", patchReq))
// //
// //	metadata := patchReq.Metadata
// //	labels := fmt.Sprintf("release=%s, tier=%s", metadata.DeploymentID, metadata.ComponentID)
// //	depErr := p.patchDeployment(labels, patchReq.Image)
// //	if depErr != nil {
// //		return nil, depErr
// //	}
// //
// //	stsErr := p.patchStatefulSet(labels, patchReq.Image)
// //	if stsErr != nil {
// //		return nil, stsErr
// //	}
// //
// //	return &provisioner.PatchDeploymentResponse{}, nil
// //}
// //
// //// Patches a Kubernetes Deployment resource.
// //func (p *KubeProvisioner) patchDeployment(labels, image string) error {
// //	logger := log.WithField("function", "patchDeployment")
// //	logger.Debug("Applying deployment patches...")
// //	logger.Debug(fmt.Sprintf("Using labels %s", labels))
// //
// //	deployments, listErr := p.deploymentsClient.List(metav1.ListOptions{
// //		LabelSelector: labels,
// //	})
// //	if listErr != nil {
// //		return listErr
// //	}
// //	logger.Debug(fmt.Sprintf("Found %d deployments", len(deployments.Items)))
// //
// //	for _, d := range deployments.Items {
// //		// Update the deployment.
// //		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
// //			d.Spec.Template.Spec.Containers[0].Image = image
// //			_, updateErr := p.deploymentsClient.Update(&d)
// //			return updateErr
// //		})
// //		if retryErr != nil {
// //			return retryErr
// //		}
// //	}
// //
// //	return nil
// //}
// //
// //// Patches a Kubernetes StatefulSet resource.
// //func (p *KubeProvisioner) patchStatefulSet(labels, image string) error {
// //	logger := log.WithField("function", "patchStatefulSet")
// //	logger.Debug("Applying statefulset patches...")
// //	logger.Debug(fmt.Sprintf("Using labels %s", labels))
// //
// //	sts, listErr := p.stsClient.List(metav1.ListOptions{
// //		LabelSelector: labels,
// //	})
// //	if listErr != nil {
// //		return listErr
// //	}
// //	logger.Debug(fmt.Sprintf("Found %d statefulsets", len(sts.Items)))
// //
// //	for _, s := range sts.Items {
// //		// Update the statefulset.
// //		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
// //			// Workers have a sidecar.
// //			s.Spec.Template.Spec.Containers[0].Image = image
// //			s.Spec.Template.Spec.Containers[1].Image = image
// //			_, updateErr := p.stsClient.Update(&s)
// //			return updateErr
// //		})
// //		if retryErr != nil {
// //			return retryErr
// //		}
// //	}
// //
// //	return nil
// //}

func BuildResult(err error) *proto.Result {
	if err != nil {
		return &proto.Result {
			Success: true,
			Message: "",
		}
	}
	return &proto.Result {
		Success: false,
		Message: err.Error(),
	}
}