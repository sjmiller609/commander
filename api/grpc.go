package api

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/astronomerio/commander/pkg/proto"
	"github.com/astronomerio/commander/provisioner"
	"github.com/astronomerio/commander/utils"
)

const (
	port = ":50051"
)

// Struct all gRPC methods will be implemented on
type GRPCServer struct{
	grpc *grpc.Server
	provisioner provisioner.Provisioner
}

// Creates a new gRPC server instance
func NewGRPC(prov provisioner.Provisioner) (*GRPCServer) {
	apiServer := GRPCServer{
		grpc: grpc.NewServer(),
		provisioner: prov,
	}

	proto.RegisterCommanderServer(apiServer.grpc, &apiServer)

	// Register reflection service on gRPC server.
	reflection.Register(apiServer.grpc)

	return &apiServer
}

// Binds a port for the gRPC server to listen to
func (s *GRPCServer) Serve(port string) error {
	logger := log.WithField("function", "Serve")
	logger.Debug("Starting gRPC server")

	port = utils.EnsurePrefix(port, ":")

	// listen to a port
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// bind port to gRPC server
	if err := s.grpc.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}