package libgrpc

import (
	"goproject/internal/handlers"
	"goproject/protos/gen/librarypb"
)

// GRPCServer реализует интерфейс LibraryServer из proto
type GRPCServer struct {
	librarypb.UnimplementedLibraryServer
	useCase handlers.UseCase
}

func NewGRPCServer(useCase handlers.UseCase) *GRPCServer {
	return &GRPCServer{useCase: useCase}
}
