package event

import (
	"time"
)

// Portfolio bir portföyü temsil eder
type Portfolio struct {
	PortID     int    `json:"portID"`
	Name       string `json:"name"`
	UserID     string `json:"userID"`
	CreatedAt  string `json:"createdAt"`
	LastUpdate string `json:"lastUpdate"`
}

// PortfolioReportPayload portfolio.report olayının payload'ını tanımlar
type PortfolioReportPayload struct {
	Portfolios []Portfolio `json:"portfolios"`
}

// NewPortfolioReportEvent yeni bir portfolio report event'i oluşturur
func NewPortfolioReportEvent(portfolios []Portfolio) (BaseEvent, error) {
	payload := PortfolioReportPayload{
		Portfolios: portfolios,
	}

	return NewBaseEvent(PortfolioReport, payload)
}

// CreateSamplePortfolios örnek portföy verileri oluşturur
func CreateSamplePortfolios() []Portfolio {
	return []Portfolio{
		{
			PortID:     1,
			Name:       "My Tech Portfolio",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -6, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     2,
			Name:       "Retirement Fund",
			UserID:     "user456",
			CreatedAt:  time.Now().AddDate(0, -4, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -5).Format("2006-01-02 15:04:05"),
		},
		{
			PortID:     3,
			Name:       "Growth Portfolio",
			UserID:     "user123",
			CreatedAt:  time.Now().AddDate(0, -2, 0).Format("2006-01-02 15:04:05"),
			LastUpdate: time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		},
	}
} 