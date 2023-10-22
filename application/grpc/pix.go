package grpc

import (
	"github.com/keyo-oliveira/codepix/application/grpc/pb"
	"github.com/keyo-oliveira/codepix/application/usecase"
	"golang.org/x/net/context"
)

type PixGrpcService struct {
	PixUseCase usecase.PixUseCase
	pb.UnimplementedPixServiceServer
}

func (p *PixGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	//key, err := p.PixUseCase.PixKeyRepository.RegisterKey(in.Key, in.Kind, in.AccountId)
	//if err != nil {
	//	return &pb.PixKeyCreatedResult{
	//			Status: "not created",
	//		Error:  err.Error(),
	//	}, err
	//}

	return &pb.PixKeyCreatedResult{
		Id:     "MOCKED KEY",
		Status: "created",
	}, nil
}

func (p *PixGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixKey, err := p.PixUseCase.PixKeyRepository.FindKeyByKind(in.Key, in.Kind)
	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id:   pixKey.ID,
		Kind: pixKey.Kind,
		Key:  pixKey.Key,
		Account: &pb.Account{
			AccountId:     pixKey.AccountID,
			AccountNumber: pixKey.Account.Number,
			BankId:        pixKey.Account.BankID,
			BankName:      pixKey.Account.Bank.Name,
			OwnerName:     pixKey.Account.OwnerName,
			CreatedAt:     pixKey.Account.CreatedAT.String(),
		},
		CreatedAt: pixKey.CreatedAT.String(),
	}, nil
}

func NewPixGrpcService(usecase usecase.PixUseCase) *PixGrpcService {
	return &PixGrpcService{PixUseCase: usecase}
}
