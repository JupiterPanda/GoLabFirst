package libgrpc

import (
	"goproject/internal/handlers"
	"goproject/protos/gen"
)

// GRPCServer реализует интерфейс LibraryServer из proto
type GRPCServer struct {
	gen.UnimplementedLibraryServer
	useCase handlers.UseCase
}

func NewGRPCServer(useCase handlers.UseCase) *GRPCServer {
	return &GRPCServer{useCase: useCase}
}
