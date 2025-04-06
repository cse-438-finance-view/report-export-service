package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/report"
)

// PortfolioReportHandler portfolio.report olaylarını işleyen yapı
type PortfolioReportHandler struct {
	PDFGenerator *report.PDFGenerator
}

// NewPortfolioReportHandler yeni bir portfolio report handler oluşturur
func NewPortfolioReportHandler() *PortfolioReportHandler {
	// PDF Generator oluştur
	pdfGenerator, err := report.NewPDFGenerator("reports")
	if err != nil {
		log.Printf("Warning: Failed to create PDF generator: %v. PDF reports will not be generated.", err)
		return &PortfolioReportHandler{PDFGenerator: nil}
	}
	
	return &PortfolioReportHandler{
		PDFGenerator: pdfGenerator,
	}
}

// EventType bu işleyicinin hangi event tipini işlediğini belirtir
func (h *PortfolioReportHandler) EventType() event.EventType {
	return event.PortfolioReport
}

// Handle portfolio.report olayını işler
func (h *PortfolioReportHandler) Handle(ctx context.Context, evt event.BaseEvent) error {
	// Payload'ı doğru tipe dönüştür
	var payload event.PortfolioReportPayload
	if err := evt.ParsePayload(&payload); err != nil {
		return fmt.Errorf("failed to parse portfolio report payload: %w", err)
	}

	// Log işlemi
	log.Printf("Processing portfolio report event with %d portfolios", len(payload.Portfolios))
	
	// Her bir portföy için ayrıntı loglar
	for _, portfolio := range payload.Portfolios {
		log.Printf("Portfolio details: ID=%d, Name=%s, UserID=%s, CreatedAt=%s, LastUpdate=%s",
			portfolio.PortID, portfolio.Name, portfolio.UserID, portfolio.CreatedAt, portfolio.LastUpdate)
	}
	
	// PDF raporu oluştur
	if h.PDFGenerator != nil {
		log.Println("Generating PDF report for portfolios...")
		
		// İşlemin biraz zaman aldığını simüle etmek için
		time.Sleep(200 * time.Millisecond)
		
		// PDF oluştur
		filePath, err := h.PDFGenerator.GeneratePortfolioReport(payload.Portfolios)
		if err != nil {
			log.Printf("Error generating PDF report: %v", err)
			// Hatayı yut ve devam et - rapor oluşturulamazsa da işlem başarılı sayılabilir
		} else {
			log.Printf("PDF report successfully generated at: %s", filePath)
		}
	} else {
		log.Println("PDF generator not available, skipping report generation")
	}
	
	// Rapor oluşturma işleminin tamamlandığını belirt
	log.Printf("Portfolio report processing completed for %d portfolios", len(payload.Portfolios))
	
	return nil
} 