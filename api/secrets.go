package api

import (
	"github.com/astronomer/commander/pkg/proto"
	"golang.org/x/net/context"
)

func (s *GRPCServer) GetSecret(ctx context.Context, in *proto.GetSecretRequest) (*proto.GetSecretResponse, error) {
	return s.provisioner.GetSecret(in)
}

func (s *GRPCServer) SetSecret(ctx context.Context, in *proto.SetSecretRequest) (*proto.SetSecretResponse, error) {
	return s.provisioner.SetSecret(in)
}
