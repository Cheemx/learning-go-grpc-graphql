package server

import (
	"context"
	"log"

	"github.com/Cheemx/learning-go-grpc-graphql/internal/entities"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/repo"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/golang_protobuf_brand"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CRUDServiceServer implements CRUDServer
type CRUDServiceServer struct {
	repo *repo.GenericRepo[entities.Brand]
}

func NewCRUDServiceServer(brandRepo *repo.GenericRepo[entities.Brand]) *CRUDServiceServer {
	return &CRUDServiceServer{
		repo: brandRepo,
	}
}

func (c CRUDServiceServer) Create(context.Context, *golang_protobuf_brand.ProtoBrandRepo_ProtoBrand) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, nil
}

func (c CRUDServiceServer) GetOne(_ context.Context, id *wrapperspb.Int64Value) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	brand, err := (*c.repo).GetOne(uint(id.Value))
	if err != nil {
		log.Printf("failed to get Brand: %v", err)
		return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, err
	}
	return repo.ToProtoBrand(brand), nil
}

func (c CRUDServiceServer) GetList(_ *emptypb.Empty, stream golang_protobuf_brand.CRUD_GetListServer) error {
	for _, brand := range (*c.repo).GetList() {
		if err := stream.Send(repo.ToProtoBrand(brand)); err != nil {
			return err
		}
	}
	return nil
}

func (c CRUDServiceServer) Update(_ context.Context, message *golang_protobuf_brand.UpdateRequest) (*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand, error) {
	brand, err := (*c.repo).Update(uint(message.ID.Value), repo.ToBrand(message.Brand))
	if err != nil {
		log.Printf("failed to update Brand: %v", err)
		return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{}, err
	}
	return repo.ToProtoBrand(brand), nil
}

func (c CRUDServiceServer) Delete(_ context.Context, message *wrapperspb.Int64Value) (*wrapperspb.BoolValue, error) {
	success, err := (*c.repo).DeleteOne(uint(message.Value))
	if err != nil {
		log.Printf("failed to delete Brand: %v", err)
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(success), nil
}
