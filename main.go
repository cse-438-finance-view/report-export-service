package main

import (
	"log"

	"github.com/burakmike/report-export-service/pkg/service"
)

func main() {
	log.Println("Starting Report Export Service...")
	svc := service.NewService()
	err := svc.Start()
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
	log.Println("Report Export Service started successfully")
	svc.WaitForSignal()
}
