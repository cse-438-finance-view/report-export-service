package handler_test

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/handler"
	"github.com/burakmike/report-export-service/pkg/report"
)

func TestIntegration_Handler_DB_and_PDF(t *testing.T) {
	// Setup a sqlmock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	// Prepare sample portfolios
	portfolios := event.CreateSamplePortfolios()

	// Expect one Exec per portfolio
	sql := "INSERT INTO reports(user_id, type) VALUES($1, $2)"
	for _, p := range portfolios {
		mock.ExpectExec(regexp.QuoteMeta(sql)).
			WithArgs(p.UserID, string(event.PortfolioReport)).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Create a PDF generator with a temporary directory
	dir := t.TempDir()
	pdfGen, err := report.NewPDFGenerator(dir)
	if err != nil {
		t.Fatalf("NewPDFGenerator error: %v", err)
	}

	// Initialize the handler
	h := handler.NewPortfolioReportHandler(db, pdfGen)

	// Create an event for portfolios
	evt, err := event.NewPortfolioReportEvent(portfolios)
	if err != nil {
		t.Fatalf("NewPortfolioReportEvent error: %v", err)
	}

	// Handle the event
	if err := h.Handle(context.Background(), evt); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	// Verify a PDF file was created
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	found := false
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".pdf" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected at least one PDF file in %s, found none", dir)
	}

	// Ensure all DB expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("DB expectations not met: %v", err)
	}
} 