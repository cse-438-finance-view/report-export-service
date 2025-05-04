package service

import (
	"os"
	"strings"
	"testing"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/handler"
	"github.com/burakmike/report-export-service/pkg/report"
)

func TestSetupHandlers(t *testing.T) {
	s := &Service{
		Registry:     handler.NewHandlerRegistry(),
		PDFGenerator: nil,
		DB:           nil,
	}
	s.SetupHandlers()
	h := s.Registry.GetHandler(event.PortfolioReport)
	if h == nil {
		t.Errorf("Expected handler for event %v, got nil", event.PortfolioReport)
	}
	// Ensure the type is correct
	if _, ok := h.(*handler.PortfolioReportHandler); !ok {
		t.Errorf("Expected *PortfolioReportHandler, got %T", h)
	}
}

func TestGeneratePortfolioReportPDF_NotInitialized(t *testing.T) {
	s := &Service{PDFGenerator: nil}
	_, err := s.GeneratePortfolioReportPDF(nil)
	if err == nil || err.Error() != "PDF generator not initialized" {
		t.Errorf("GeneratePortfolioReportPDF(nil) error = %v; want \"PDF generator not initialized\"", err)
	}
}

func TestGeneratePortfolioReportPDF_Success(t *testing.T) {
	dir := t.TempDir()
	gen, err := report.NewPDFGenerator(dir)
	if err != nil {
		t.Fatalf("NewPDFGenerator error: %v", err)
	}
	s := &Service{PDFGenerator: gen}
	portfolios := []event.Portfolio{
		{PortID: 1, Name: "Test", UserID: "user1", CreatedAt: "2023-01-01 00:00:00", LastUpdate: "2023-01-02 00:00:00"},
	}
	path, err := s.GeneratePortfolioReportPDF(portfolios)
	if err != nil {
		t.Fatalf("GeneratePortfolioReportPDF error: %v", err)
	}
	// Path ends with .pdf
	if !strings.HasSuffix(path, ".pdf") {
		t.Errorf("Expected path to end with .pdf, got %s", path)
	}
	// File exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("PDF file not created at %s", path)
	}
} 