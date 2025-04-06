package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/burakmike/report-export-service/pkg/report"
)

func main() {
	// Komut satırı argümanlarını kontrol et
	outputDir := "reports"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	fmt.Printf("PDF rapor oluşturma testi başlıyor...\n")
	fmt.Printf("Raporlar şu dizine kaydedilecek: %s\n\n", outputDir)

	// PDF oluşturucu oluştur
	pdfGenerator, err := report.NewPDFGenerator(outputDir)
	if err != nil {
		log.Fatalf("PDF generator oluşturulamadı: %v", err)
	}

	// Örnek portföy verileri oluştur
	portfolios := generateSamplePortfolios()
	fmt.Printf("%d adet örnek portföy oluşturuldu.\n", len(portfolios))

	// PDF raporu oluştur
	fmt.Println("PDF raporu oluşturuluyor...")
	filePath, err := pdfGenerator.GeneratePortfolioReport(portfolios)
	if err != nil {
		log.Fatalf("PDF raporu oluşturulurken hata: %v", err)
	}

	fmt.Printf("\nPDF raporu başarıyla oluşturuldu!\n")
	fmt.Printf("Oluşturulan dosya: %s\n", filePath)
	fmt.Printf("Tam dosya yolu: %s\n", getAbsolutePath(filePath))
}

// generateSamplePortfolios örnek portföy verileri oluşturur
func generateSamplePortfolios() []event.Portfolio {
	// Türkçe isimlerle örnek portföyler
	return []event.Portfolio{
		{
			PortID:     1,
			Name:       "Teknoloji Portföyü",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -6, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     2,
			Name:       "Emeklilik Fonu",
			UserID:     "user456",
			CreatedAt:  time.Now().AddDate(0, -4, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -5).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     3,
			Name:       "Büyüme Portföyü",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -2, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     4,
			Name:       "Gayrimenkul Yatırım Fonu",
			UserID:     "user789",
			CreatedAt:  time.Now().AddDate(0, -8, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -15).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     5,
			Name:       "Uzun Vadeli Tahvil Portföyü",
			UserID:     "user456",
			CreatedAt:  time.Now().AddDate(0, -3, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04:05"),
		},
	}
}

// getAbsolutePath dosyanın tam yolunu döndürür
func getAbsolutePath(filePath string) string {
	absPath, err := os.Getwd()
	if err != nil {
		return filePath
	}
	return fmt.Sprintf("%s/%s", absPath, filePath)
} 