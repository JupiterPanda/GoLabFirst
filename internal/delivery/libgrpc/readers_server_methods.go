package libgrpc

import (
	"context"
	"goproject/internal/models"
	"goproject/protos/gen/librarypb"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) GetAllReaders(ctx context.Context, _ *emptypb.Empty) (*librarypb.ReadersResponse, error) {
	readers, err := s.useCase.GetAllReaders(ctx)
	if err != nil {
		return nil, err
	}

	res := &librarypb.ReadersResponse{}
	for _, reader := range readers {
		res.Readers = append(res.Readers, &librarypb.Reader{
			ID:          int32(reader.ID),
			Name:        reader.Name,
			PhoneNumber: reader.PhoneNumber,
			Address:     reader.Address,
			DateOfBirth: timestamppb.New(reader.DateOfBirth),
		})
	}

	return res, nil
}

func (s *GRPCServer) GetReaderIdByName(ctx context.Context, req *librarypb.ReaderNameRequest) (*librarypb.IdResponse, error) {
	readerID, err := s.useCase.GetReaderIdByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	res := &librarypb.IdResponse{Id: int32(readerID)}
	return res, nil
}

func (s *GRPCServer) CreateReader(ctx context.Context, req *librarypb.CreateReaderRequest) (*librarypb.MessageResponse, error) {
	reader := models.Reader{
		ID:          int(req.Reader.ID),
		Name:        req.Reader.Name,
		PhoneNumber: req.Reader.PhoneNumber,
		Address:     req.Reader.Address,
		DateOfBirth: req.Reader.DateOfBirth.AsTime(),
	}
	err := s.useCase.CreateReader(ctx, reader)
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Reader created")}
	return res, nil
}

func (s *GRPCServer) DeleteReader(ctx context.Context, req *librarypb.ReaderIdRequest) (*librarypb.MessageResponse, error) {
	err := s.useCase.DeleteReader(ctx, int(req.ReaderId))
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Reader deleted")}
	return res, nil
}

func (s *GRPCServer) UpdateReaderContactInfo(ctx context.Context, req *librarypb.Reader) (*librarypb.MessageResponse, error) {

	err := s.useCase.UpdateReaderContactInfo(ctx, int(req.ID), req.PhoneNumber, req.Address)
	if err != nil {
		return nil, err
	}
	res := &librarypb.MessageResponse{Message: ("Reader info updated")}
	return res, nil
}
