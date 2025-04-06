package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/burakmike/report-export-service/pkg/event"
	"github.com/jung-kurt/gofpdf"
)

// PDFGenerator PDF rapor oluşturmak için kullanılan yapı
type PDFGenerator struct {
	OutputDir  string // Raporların kaydedileceği dizin
	ReportLogo string // Rapor logosu (opsiyonel)
}

// ReportOptions rapor oluşturma seçeneklerini belirtir
type ReportOptions struct {
	Title    string
	Subtitle string
	Logo     string
}

// NewPDFGenerator yeni bir PDF generator oluşturur
func NewPDFGenerator(outputDir string) (*PDFGenerator, error) {
	// Dizinin var olduğunu kontrol et, yoksa oluştur
	if outputDir == "" {
		outputDir = "reports"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &PDFGenerator{
		OutputDir: outputDir,
	}, nil
}

// GeneratePortfolioReport portföy verilerinden PDF raporu oluşturur
func (g *PDFGenerator) GeneratePortfolioReport(portfolios []event.Portfolio) (string, error) {
	// Varsayılan rapor seçenekleri
	options := ReportOptions{
		Title:    "Portfolio Report",
		Subtitle: fmt.Sprintf("Generated on %s", time.Now().Format("January 2, 2006")),
	}

	return g.generateReport(portfolios, options)
}

// generateReport belirtilen seçeneklerle PDF raporu oluşturur
func (g *PDFGenerator) generateReport(portfolios []event.Portfolio, options ReportOptions) (string, error) {
	// PDF dosyasını oluştur - Yatay A4 kağıdı
	pdf := gofpdf.New("L", "mm", "A4", "")
	
	// Sayfa numaralarını ekle
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	
	// Yeni sayfa ekle
	pdf.AddPage()

	// Üst bilgi - başlık ve tarih
	g.addHeader(pdf, options)
	
	// PDF içeriğini oluştur
	g.addPortfolioTable(pdf, portfolios)
	
	// Alt bilgi - copyright ve diğer bilgiler
	g.addFooter(pdf)
	
	// Dosya adını oluştur
	fileName := fmt.Sprintf("portfolio_report_%s.pdf", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(g.OutputDir, fileName)
	
	// PDF dosyasını kaydet
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save PDF file: %w", err)
	}
	
	return filePath, nil
}

// addHeader PDF'e başlık ekler
func (g *PDFGenerator) addHeader(pdf *gofpdf.Fpdf, options ReportOptions) {
	// Başlık için font ve renk ayarları
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 51, 102) // Koyu mavi
	
	// Başlık
	pdf.Cell(0, 10, options.Title)
	pdf.Ln(10)
	
	// Alt başlık
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(120, 120, 120) // Gri
	pdf.Cell(0, 10, options.Subtitle)
	pdf.Ln(5)
	
	// Ayraç çizgisi
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY()+5, 287, pdf.GetY()+5)
	pdf.Ln(15)
	
	// Başlık ve içerik arasında boşluk
	pdf.SetTextColor(0, 0, 0) // Siyah
}

// addPortfolioTable PDF'e portföy tablosunu ekler
func (g *PDFGenerator) addPortfolioTable(pdf *gofpdf.Fpdf, portfolios []event.Portfolio) {
	// Tablo başlıkları için font ayarla
	pdf.SetFont("Arial", "B", 11)
	
	// Tablo başlık renkleri
	pdf.SetFillColor(66, 133, 244) // Google mavi
	pdf.SetTextColor(255, 255, 255) // Beyaz
	pdf.SetDrawColor(66, 133, 244) // Google mavi
	
	// Sütun genişlikleri
	colWidths := []float64{20, 90, 40, 60, 60}
	
	// Tablo başlıklarını ekle
	header := []string{"ID", "Portfolio Name", "User ID", "Created", "Last Updated"}
	for i, heading := range header {
		pdf.CellFormat(colWidths[i], 8, heading, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	
	// Tablo içeriği için font ve renk ayarla
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0) // Siyah
	
	// İçerik için alternatif satır renkleri
	evenRowColor := []int{240, 240, 240} // Açık gri
	oddRowColor := []int{255, 255, 255}  // Beyaz
	
	// Tablo hücre sınır rengi
	pdf.SetDrawColor(200, 200, 200) // Açık gri
	
	// Her bir portfolyo satırını ekle
	for i, portfolio := range portfolios {
		// Alternatif satır renkleri
		if i%2 == 0 {
			pdf.SetFillColor(evenRowColor[0], evenRowColor[1], evenRowColor[2])
		} else {
			pdf.SetFillColor(oddRowColor[0], oddRowColor[1], oddRowColor[2])
		}
		
		// Portföy ID
		pdf.CellFormat(colWidths[0], 8, fmt.Sprintf("%d", portfolio.PortID), "1", 0, "C", true, 0, "")
		
		// Portföy adı
		pdf.CellFormat(colWidths[1], 8, portfolio.Name, "1", 0, "L", true, 0, "")
		
		// Kullanıcı ID
		pdf.CellFormat(colWidths[2], 8, portfolio.UserID, "1", 0, "C", true, 0, "")
		
		// Oluşturulma tarihi
		pdf.CellFormat(colWidths[3], 8, portfolio.CreatedAt, "1", 0, "C", true, 0, "")
		
		// Son güncelleme tarihi
		pdf.CellFormat(colWidths[4], 8, portfolio.LastUpdate, "1", 0, "C", true, 0, "")
		
		pdf.Ln(-1)
	}
	
	// Toplam bilgisi
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 8, fmt.Sprintf("Total Portfolios: %d", len(portfolios)), "", 0, "R", false, 0, "")
	pdf.Ln(15)
}

// addFooter PDF'e alt bilgi ekler
func (g *PDFGenerator) addFooter(pdf *gofpdf.Fpdf) {
	// Alt bilgi renk ve font ayarları
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	
	// Oluşturulma bilgisi
	generatedMsg := fmt.Sprintf("This report was automatically generated on %s", 
		time.Now().Format("January 2, 2006 at 15:04:05"))
	pdf.CellFormat(0, 10, generatedMsg, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	
	// Yasal uyarı/copyright
	pdf.CellFormat(0, 10, "© Portfolio Report Service - Confidential Information", "", 0, "L", false, 0, "")
} 