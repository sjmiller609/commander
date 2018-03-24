package api

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) ListDeployments(ctx context.Context, in *ListDeploymentsRequest) (*ListDeploymentsResponse, error) {
	fmt.Println("ListDeployments called")
	return &ListDeploymentsResponse{}, nil
}

func (s *Server) CreateDeployment(ctx context.Context, in *CreateDeploymentRequest) (*CreateDeploymentResponse, error) {
	fmt.Println("CreateDeployment called")
	return &CreateDeploymentResponse{}, nil
}

func (s *Server) UpdateDeployment(ctx context.Context, in *UpdateDeploymentRequest) (*UpdateDeploymentResponse, error) {
	fmt.Println("UpdateDeployment called")
	return &UpdateDeploymentResponse{}, nil
}

func (s *Server) UpgradeDeployment(ctx context.Context, in *UpgradeDeploymentRequest) (*UpgradeDeploymentResponse, error) {
	fmt.Println("UpgradeDeployment called")
	return &UpgradeDeploymentResponse{}, nil
}

func (s *Server) DeleteDeployment(ctx context.Context, in *DeleteDeploymentRequest) (*DeleteDeploymentResponse, error) {
	fmt.Println("DeleteDeployment called")
	return &DeleteDeploymentResponse{}, nil
}

