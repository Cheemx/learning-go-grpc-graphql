package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Cheemx/learning-go-grpc-graphql/internal/config"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/controller"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/entities"
	"github.com/Cheemx/learning-go-grpc-graphql/internal/repo"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	RegisterProductRoutes(router)
	RegisterBrandRoutes(router)

	log.Printf("Starting Server on port %s\n", cfg.Port)
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
