package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/burakmike/report-export-service/pkg/config"
	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/handler"
	"github.com/burakmike/report-export-service/pkg/rabbitmq"
	"github.com/burakmike/report-export-service/pkg/report"
	"github.com/streadway/amqp"
)

const (
	maxReconnectRetries = 5
	defaultReportDir    = "reports"
)

// Service ana uygulama yapısını temsil eder, tüm bileşenleri koordine eder
type Service struct {
	Config       config.Config
	RabbitMQ     *rabbitmq.RabbitMQClient
	Registry     *handler.HandlerRegistry
	Context      context.Context
	CancelFunc   context.CancelFunc
	PDFGenerator *report.PDFGenerator
}

// NewService yeni bir hizmet oluşturur ve bağımlılıkları kurar
func NewService() *Service {
	// Context oluştur
	ctx, cancel := context.WithCancel(context.Background())
	
	// Config yükle
	cfg := config.LoadConfigFromEnv()
	
	// Handler kayıt sistemini oluştur
	registry := handler.NewHandlerRegistry()
	
	// PDF Generator oluştur
	pdfGenerator, err := report.NewPDFGenerator(defaultReportDir)
	if err != nil {
		log.Printf("Warning: Failed to initialize PDF generator: %v. PDF reports will not be generated.", err)
		pdfGenerator = nil
	} else {
		log.Printf("PDF generator initialized. Reports will be saved to: %s", defaultReportDir)
	}
	
	// RabbitMQ client'ını oluştur
	rabbitClient := rabbitmq.NewRabbitMQClient(cfg, registry)
	
	return &Service{
		Config:       cfg,
		RabbitMQ:     rabbitClient,
		Registry:     registry,
		Context:      ctx,
		CancelFunc:   cancel,
		PDFGenerator: pdfGenerator,
	}
}

// SetupHandlers tüm event işleyicileri kaydeder
func (s *Service) SetupHandlers() {
	// Portfolio rapor işleyicisi
	portfolioHandler := handler.NewPortfolioReportHandler()
	s.Registry.RegisterHandler(portfolioHandler)
	
	// İleride başka işleyiciler de buraya eklenebilir
	
	log.Println("Event handlers registered")
}

// Start hizmeti başlatır
func (s *Service) Start() error {
	log.Println("Starting service...")

	// İşleyicileri kur
	s.SetupHandlers()

	// Raporlar dizinini oluştur (yoksa)
	if _, err := os.Stat(defaultReportDir); os.IsNotExist(err) {
		log.Printf("Creating reports directory: %s", defaultReportDir)
		if err := os.MkdirAll(defaultReportDir, 0755); err != nil {
			log.Printf("Warning: Failed to create reports directory: %v", err)
		}
	}
	
	// RabbitMQ'ya bağlan
	err := s.RabbitMQ.Connect()
	if err != nil {
		return err
	}
	
	// Exchange ve kuyrukları oluştur
	err = s.RabbitMQ.SetupExchangeAndQueues()
	if err != nil {
		s.RabbitMQ.Close()
		return err
	}
	
	// Mesajları tüketmeye başla
	err = s.RabbitMQ.ConsumeMessages(s.Context)
	if err != nil {
		s.RabbitMQ.Close()
		return err
	}
	
	// Bağlantı kapanışını izleyen kanal
	closeChan := make(chan *amqp.Error)
	s.RabbitMQ.Connection.NotifyClose(closeChan)
	
	// Kapanış sinyali için bir dinleyici başlat
	go s.monitorConnection(closeChan)
	
	log.Println("Service started successfully")
	
	return nil
}

// monitorConnection RabbitMQ bağlantısını izler, kapanırsa yeniden bağlanmayı dener
func (s *Service) monitorConnection(closeChan chan *amqp.Error) {
	for {
		select {
		case err := <-closeChan:
			if err != nil {
				log.Printf("RabbitMQ connection closed: %v", err)
				
				// Yeniden bağlanmayı dene
				reconnectErr := s.RabbitMQ.Reconnect(s.Context, maxReconnectRetries)
				if reconnectErr != nil {
					log.Fatalf("Failed to reconnect to RabbitMQ: %v", reconnectErr)
				}
				
				// Yeniden bağlantı kapanış bildirimi al
				closeChan = make(chan *amqp.Error)
				s.RabbitMQ.Connection.NotifyClose(closeChan)
			}
		case <-s.Context.Done():
			log.Println("Connection monitor shutting down")
			return
		}
	}
}

// Stop hizmeti durdurur ve kaynakları temizler
func (s *Service) Stop() {
	log.Println("Shutting down service...")
	
	// Context'i iptal et
	s.CancelFunc()
	
	// RabbitMQ bağlantısını kapat
	if s.RabbitMQ != nil {
		s.RabbitMQ.Close()
	}
	
	log.Println("Service shutdown complete")
}

// WaitForSignal servisin sonlanması için sinyal bekler
func (s *Service) WaitForSignal() {
	// Sinyal dinleme kanalı
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Sinyal gelene kadar bekle
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
	
	// Temiz bir kapatma başlat
	s.Stop()
}

// PublishPortfolioReport portföy rapor olayı yayınlar
func (s *Service) PublishPortfolioReport(portfolios []event.Portfolio) error {
	return s.RabbitMQ.PublishPortfolioReport(portfolios)
}

// GeneratePortfolioReportPDF portföy verilerinden PDF raporu oluşturur
func (s *Service) GeneratePortfolioReportPDF(portfolios []event.Portfolio) (string, error) {
	if s.PDFGenerator == nil {
		return "", fmt.Errorf("PDF generator not initialized")
	}
	
	return s.PDFGenerator.GeneratePortfolioReport(portfolios)
} 