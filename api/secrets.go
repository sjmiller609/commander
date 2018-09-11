package api

import (
	"github.com/astronomerio/commander/pkg/proto"
	"golang.org/x/net/context"
)

func (s *GRPCServer) GetSecret(ctx context.Context, in *proto.GetSecretRequest) (*proto.GetSecretResponse, error) {
	return s.provisioner.GetSecret(in)
}
