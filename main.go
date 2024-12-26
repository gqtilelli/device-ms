package main

import (
	"context"
	"log"
	"net/http"

	"github.com/device-ms/controller"
	"github.com/device-ms/handler"
	"github.com/device-ms/mongo"
)

func main() {
	ctx := context.Background()
	r := initRouter(ctx)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func initRouter(ctx context.Context) (router handler.Router) {
	deviceRepository, err := mongo.CreateDeviceRepo(ctx)
	if err != nil {
		log.Fatal("Could not initialize device repository: " + err.Error())
	}

	service := controller.New(ctx, deviceRepository)

	return handler.NewDeviceRouter(service)
}
