package provisioner

import (
	"github.com/astronomerio/commander/pkg/proto"
)
// Provisioner types are capable of scheduling and various maintenace tasks
// of containers running on various container orchestrators.
type Provisioner interface {
	// ListDeployments(organizationID string) (*ListDeploymentResponse, error)
	InstallDeployment(request *proto.CreateDeploymentRequest) (*proto.CreateDeploymentResponse, error)
	UpdateDeployment(request *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse, error)
	// UpgradeDeployment()
	DeleteDeployment(request *proto.DeleteDeploymentRequest) (*proto.DeleteDeploymentResponse, error)
	// FetchDeployments()
	// PatchDeployment(patchReq *PatchDeploymentRequest) (*PatchDeploymentResponse, error)
}

// ListDeploymentResponse is a response from listing deployments.
type ListDeploymentResponse struct {
	Items []string
}

// DeploymentMetadata uniquely describes a deploymnet.
type DeploymentMetadata struct {
	DeploymentID string `json:"deploymentId,omitempty"`
	ComponentID  string `json:"componentId,omitempty"`
}

// PatchDeploymentRequest is a request to patch an existing deployment.
type PatchDeploymentRequest struct {
	Metadata DeploymentMetadata
	Image    string `json:"image,omitempty"`
	// Workers int    `json:"workers,omitempty"`
}

// PatchDeploymentResponse is a response from patching an existing deployment.
type PatchDeploymentResponse struct {

}
