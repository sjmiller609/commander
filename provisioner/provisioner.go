package provisioner

// Provisioner types are capable of scheduling and various maintenace tasks
// of containers running on various container orchestrators.
type Provisioner interface {
	ListDeployments(organizationID string) (*ListDeploymentResponse, error)
	PatchDeployment(deploymentID string, patchReq *PatchDeploymentRequest) (*PatchDeploymentResponse, error)
}

// ListDeploymentResponse is a response from listing deployments.
type ListDeploymentResponse struct {
	Items []string
}

// PatchDeploymentRequest is a request to patch an existing deployment.
type PatchDeploymentRequest struct {
	Image string `json:"image,omitempty"`
	// Workers int    `json:"workers,omitempty"`
}

// PatchDeploymentResponse is a response from patching an existing deployment.
type PatchDeploymentResponse struct {
}
