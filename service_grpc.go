package lifetime

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

// NewGRPCService returns a service that will run listen and serve the given
// GRPC server.
func NewGRPCService(server *grpc.Server, listenAddress string) Service {
	return &grpcService{
		server:        server,
		listenAddress: listenAddress,
	}
}

// grpcService is an implementation of Service that will listen and serve the given
// HTTP server.
type grpcService struct {
	server        *grpc.Server
	listenAddress string
}

// Start will start the service.
// This is a blocking call and should block for the lifetime of the service.
// Returns an error which is treated as fatal.
func (service *grpcService) Start() error {
	lis, err := net.Listen("tcp", service.listenAddress)
	if err != nil {
		return fmt.Errorf("could not listen on tcp address: %w", err)
	}
	err = service.server.Serve(lis)
	if err == nil {
		return nil
	}
	// ErrServerStopped is returned when we call server.Close() from Service.Stop
	// so we shouldn't treat it as a breaking error.
	if err == grpc.ErrServerStopped {
		return nil
	}
	return err
}

// Stop will stop the service.
// Stop is not called if Start returned an error.
func (service *grpcService) Stop() {
	service.server.GracefulStop()
}
