package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheemx/learning-go-grpc-graphql/internal/config"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/entities"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/repo"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/golang_protobuf_brand"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var brandRepo repo.GenericRepo[entities.Brand] = repo.NewBrandRepo()

	log.Printf("Starting Server on port %s\n", cfg.Port)
	go StartRPCServer(&brandRepo)

	go StartRPCGatewayServer()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	log.Println("Termination signal received. Exiting...")
}

func StartRPCGatewayServer() {
	gwmux := runtime.NewServeMux()
	err := golang_protobuf_brand.RegisterCRUDHandlerFromEndpoint(context.Background(), gwmux, ":6002", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    ":6000",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway pn http://localhost:6002")
	log.Fatalln(gwServer.ListenAndServe())
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
