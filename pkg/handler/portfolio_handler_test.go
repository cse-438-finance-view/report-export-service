package handler

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/report"
)

func TestPortfolioReportHandler_InvalidPayload(t *testing.T) {
	h := NewPortfolioReportHandler(nil, nil)
	// Invalid JSON payload
	evt := event.BaseEvent{EventType: event.PortfolioReport, Timestamp: "x", Payload: json.RawMessage("invalid")}
	err := h.Handle(context.Background(), evt)
	if err == nil {
		t.Fatal("Expected error for invalid payload, got nil")
	}
}

func TestPortfolioReportHandler_NoPDFNoDB(t *testing.T) {
	// Valid payload, but no PDFGenerator and no DB
	payload := event.PortfolioReportPayload{Portfolios: event.CreateSamplePortfolios()}
	raw, _ := json.Marshal(payload)
	evt := event.BaseEvent{EventType: event.PortfolioReport, Timestamp: "x", Payload: raw}
	h := NewPortfolioReportHandler(nil, nil)
	err := h.Handle(context.Background(), evt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPortfolioReportHandler_WithPDF(t *testing.T) {
	// Valid payload with PDF generator
	dir := t.TempDir()
	pdfGen, err := report.NewPDFGenerator(dir)
	if err != nil {
		t.Fatalf("NewPDFGenerator error: %v", err)
	}
	payload := event.PortfolioReportPayload{Portfolios: []event.Portfolio{
		{PortID: 1, Name: "n", UserID: "u", CreatedAt: "c", LastUpdate: "l"},
	}}
	raw, _ := json.Marshal(payload)
	evt := event.BaseEvent{EventType: event.PortfolioReport, Timestamp: "x", Payload: raw}
	h := NewPortfolioReportHandler(nil, pdfGen)
	err = h.Handle(context.Background(), evt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Check PDF file exists in output directory
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	found := false
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".pdf") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected a PDF file in %s, found none", dir)
	}
} 