package main

import (
	mongosh "device-control/Mongosh"
	"device-control/config"
	"device-control/internal/service"
	"device-control/internal/storage"
	"device-control/pkg"
	"log"
)

func main() {
	cfg := config.LOAD("./")
	db, err := mongosh.InitMongo(cfg)
	if err != nil {
		log.Fatal(err)
	}

	deviceRepo := storage.NewDeviceControlRepo(*db)
	service := service.NewService(deviceRepo)

	copyService := pkg.NewCopyService(*service)

	log.Printf("DeviceControl service running on :%s port", cfg.DeviceControlPort)
	if err := copyService.RUN(cfg); err != nil {
		log.Fatal(err)
	}
}
