package libgrpc

import (
	"context"
	"goproject/internal/models"
	"goproject/protos/gen"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) GetAllReaders(ctx context.Context, _ *emptypb.Empty) (*gen.ReadersResponse, error) {
	readers, err := s.useCase.GetAllReaders(ctx)
	if err != nil {
		return nil, err
	}

	res := &gen.ReadersResponse{}
	for _, reader := range readers {
		res.Readers = append(res.Readers, &gen.Reader{
			ID:          int32(reader.ID),
			Name:        reader.Name,
			PhoneNumber: reader.PhoneNumber,
			Address:     reader.Address,
			DateOfBirth: timestamppb.New(reader.DateOfBirth),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetReaderIdByName(ctx context.Context, req *gen.ReaderNameRequest) (*gen.IdResponse, error) {
	readerID, err := s.useCase.GetReaderIdByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	res := &gen.IdResponse{Id: int32(readerID)}
	return res, nil
}

func (s *GRPCServer) CreateReader(ctx context.Context, req *gen.CreateReaderRequest) (*gen.MessageResponse, error) {
	reader := models.Reader{
		Name:        req.Reader.Name,
		PhoneNumber: req.Reader.PhoneNumber,
		Address:     req.Reader.Address,
		DateOfBirth: req.Reader.DateOfBirth.AsTime(),
	}
	err := s.useCase.CreateReader(ctx, reader)
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Reader created"}
	return res, nil
}

func (s *GRPCServer) DeleteReader(ctx context.Context, req *gen.ReaderIdRequest) (*gen.MessageResponse, error) {
	err := s.useCase.DeleteReader(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Reader deleted"}
	return res, nil
}

func (s *GRPCServer) UpdateReaderContactInfo(ctx context.Context, req *gen.Reader) (*gen.MessageResponse, error) {

	err := s.useCase.UpdateReaderContactInfo(ctx, int(req.ID), req.PhoneNumber, req.Address)
	if err != nil {
		return nil, err
	}
	res := &gen.MessageResponse{Message: "Reader info updated"}
	return res, nil
}
