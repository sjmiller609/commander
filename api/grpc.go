package api

import (
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/astronomerio/commander/provisioner"
)

const (
	port = ":50051"
)

// Struct all gRPC methods will be implemented on
type Server struct{
	grpc *grpc.Server
	provisioner *provisioner.Provisioner
}

// Creates a new gRPC server instance
func NewServer() (*Server) {
	apiServer := Server{
		grpc: grpc.NewServer(),
	}
	RegisterCommanderServer(apiServer.grpc, &apiServer)

	// Register reflection service on gRPC server.
	reflection.Register(apiServer.grpc)

	return &apiServer
}

// Binds a port for the gRPC server to listen to
func (s *Server) Serve(port string) error {
	// ensure port has a colon prefix
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

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