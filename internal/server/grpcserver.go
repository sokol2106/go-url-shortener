package server

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers"
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers/grpcservice/proto"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"google.golang.org/grpc"
	"net"
)

// grpcServer структура gRPC-сервера.
type grpcServer struct {
	server *grpc.Server
}

// NewGRPCServer создает gRPC-сервер.
func NewGRPCServer() *grpcServer {
	return &grpcServer{
		server: grpc.NewServer(),
	}
}

// StartGRPCServer создает и запускает gRPC-сервер.
func (g *grpcServer) StartGRPCServer(address string, srvSh *service.ShortURL, srvAu *service.Authorization, subnet string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpchandlers.NewURLShortenerServer(srvSh, srvAu, subnet)
	proto.RegisterURLShortenerServer(g.server, server)

	return g.server.Serve(lis)
}

// StopGRPCServer останавливает gRPC-сервер.
func (g *grpcServer) StopGRPCServer() {
	g.server.GracefulStop()
}
