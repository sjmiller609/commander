package api

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/astronomerio/commander/pkg/proto"
)

func (s *GRPCServer) FetchDeployment(ctx context.Context, in *proto.FetchDeploymentRequest) (*proto.FetchDeploymentResponse, error) {
	return &proto.FetchDeploymentResponse{}, nil
}

func (s *GRPCServer) CreateDeployment(ctx context.Context, in *proto.CreateDeploymentRequest) (*proto.CreateDeploymentResponse, error) {
	fmt.Printf("+%v\n", in)
	response := &proto.CreateDeploymentResponse{}
	//response := s.provisioner.InstallDeployment(in)
	fmt.Println("CreateDeployment called")
	return response, nil
}

func (s *GRPCServer) UpdateDeployment(ctx context.Context, in *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse, error) {
	response := s.provisioner.UpdateDeployment(in)
	fmt.Println("CreateDeployment called")
	return response, nil
}

func (s *GRPCServer) UpgradeDeployment(ctx context.Context, in *proto.UpgradeDeploymentRequest) (*proto.UpgradeDeploymentResponse, error) {
	fmt.Println("UpgradeDeployment called")
	return &proto.UpgradeDeploymentResponse{}, nil
}

func (s *GRPCServer) DeleteDeployment(ctx context.Context, in *proto.DeleteDeploymentRequest) (*proto.DeleteDeploymentResponse, error) {
	fmt.Println("DeleteDeployment called")
	return &proto.DeleteDeploymentResponse{}, nil
}