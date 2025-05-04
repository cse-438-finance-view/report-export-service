package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/burakmike/report-export-service/pkg/event"
)

func TestGeneratePortfolioReport(t *testing.T) {
	// Use a temporary directory for output
	dir := t.TempDir()
	gen, err := NewPDFGenerator(dir)
	if err != nil {
		t.Fatalf("NewPDFGenerator error: %v", err)
	}

	portfolios := []event.Portfolio{
		{PortID: 1, Name: "Test1", UserID: "user1", CreatedAt: "2023-01-01 00:00:00", LastUpdate: "2023-01-02 00:00:00"},
		{PortID: 2, Name: "Test2", UserID: "user2", CreatedAt: "2023-02-01 00:00:00", LastUpdate: "2023-02-02 00:00:00"},
	}

	filePath, err := gen.GeneratePortfolioReport(portfolios)
	if err != nil {
		t.Fatalf("GeneratePortfolioReport error: %v", err)
	}

	// File should exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Expected report file %s to exist, but not found", filePath)
	}

	// File name should end with .pdf
	if !strings.HasSuffix(filePath, ".pdf") {
		t.Errorf("Expected file to have .pdf extension, got %s", filePath)
	}

	// File should be inside the output directory
	rel, err := filepath.Rel(dir, filePath)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}
	if strings.HasPrefix(rel, "..") {
		t.Errorf("Generated file path %s is not inside temp dir %s", filePath, dir)
	}

	// File size should be greater than zero
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("Expected file size to be > 0, got 0")
	}
} 