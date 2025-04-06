package main

import (
	"log"

	"github.com/burakmike/report-export-service/pkg/service"
)

func main() {
	log.Println("Starting Report Export Service...")
	
	// Yeni hizmet oluştur
	svc := service.NewService()
	
	// Hizmeti başlat
	err := svc.Start()
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
	
	log.Println("Report Export Service started successfully")
	
	// Kapatma sinyali bekle
	svc.WaitForSignal()
}
