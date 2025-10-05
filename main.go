package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Cheemx/learning-go-grpc-graphql/internal/config"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/controller"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/entities"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/repo"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/golang_protobuf_brand"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/server"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var brandRepo repo.GenericRepo[entities.Brand] = repo.NewBrandRepo()

	router := mux.NewRouter()
	RegisterProductRoutes(router)
	RegisterBrandRoutes(router)

	log.Printf("Starting Server on port %s\n", cfg.Port)
	go StartRPCServer(&brandRepo)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func RegisterProductRoutes(router *mux.Router) {
	var muxBase = "/api/products"
	router.HandleFunc(muxBase, controller.GetProducts).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controller.GetProductById).Methods("GET")
	router.HandleFunc(muxBase, controller.CreateProduct).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controller.UpdateProduct).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controller.DeleteProduct).Methods("DELETE")
}

func RegisterBrandRoutes(router *mux.Router) {
	var brandRepo repo.GenericRepo[entities.Brand] = repo.NewBrandRepo()
	NewGenericRouter[entities.Brand, *repo.BrandRepo]("/api/brands", router, &brandRepo)
}

func StartRPCServer(brandRepo *repo.GenericRepo[entities.Brand]) {
	lis, err := net.Listen("tcp", ":6002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	golang_protobuf_brand.RegisterCRUDServer(s, server.NewCRUDServiceServer(brandRepo))

	log.Printf("gRPC server listening on port %v\n", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
