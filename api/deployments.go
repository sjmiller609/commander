package api

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/astronomerio/commander/pkg/proto"
)

func (s *GRPCServer) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Received: time.Now().UnixNano() / int64(time.Millisecond),
	}, nil
}

func (s *GRPCServer) GetDeployment(ctx context.Context, in *proto.GetDeploymentRequest) (*proto.GetDeploymentResponse, error) {
	return &proto.GetDeploymentResponse{}, nil
}

func (s *GRPCServer) CreateDeployment(ctx context.Context, in *proto.CreateDeploymentRequest) (*proto.CreateDeploymentResponse, error) {
	response, err := s.provisioner.InstallDeployment(in)
	return response, err
}

func (s *GRPCServer) UpdateNamespace(ctx context.Context, in *proto.UpdateNamespaceRequest) (*proto.UpdateNamespaceResponse, error) {
	response, err := s.provisioner.UpdateNamespace(in)
	return response, err
}

func (s *GRPCServer) UpdateDeployment(ctx context.Context, in *proto.UpdateDeploymentRequest) (*proto.UpdateDeploymentResponse, error) {
	response, err := s.provisioner.UpdateDeployment(in)
	return response, err
}

func (s *GRPCServer) UpgradeDeployment(ctx context.Context, in *proto.UpgradeDeploymentRequest) (*proto.UpgradeDeploymentResponse, error) {
	response, err := s.provisioner.UpgradeDeployment(in)
	return response, err
}

func (s *GRPCServer) DeleteDeployment(ctx context.Context, in *proto.DeleteDeploymentRequest) (*proto.DeleteDeploymentResponse, error) {
	fmt.Println("DeleteDeployment called")
	response, err := s.provisioner.DeleteDeployment(in)
	return response, err
}
