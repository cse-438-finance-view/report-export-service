package event

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventType tüm event tiplerinin tanımlandığı tip
type EventType string

// Event tipleri burada tanımlanır
const (
	PortfolioReport EventType = "portfolio.report"
	// Diğer event tipleri buraya eklenebilir
	// Örnek: UserRegistered EventType = "user.registered"
	// Örnek: OrderCreated EventType = "order.created"
)

// BaseEvent tüm eventlerin içermesi gereken temel alanları tanımlar
type BaseEvent struct {
	EventType EventType       `json:"event_type"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// NewBaseEvent yeni bir temel event oluşturur
func NewBaseEvent(eventType EventType, payload interface{}) (BaseEvent, error) {
	// Payload'ı JSON formatına dönüştür
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return BaseEvent{}, fmt.Errorf("payload serialization error: %w", err)
	}

	// Temel event oluştur
	return BaseEvent{
		EventType: eventType,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Payload:   payloadBytes,
	}, nil
}

// ParseEvent bir JSON mesajını BaseEvent yapısına dönüştürür
func ParseEvent(data []byte) (BaseEvent, error) {
	var event BaseEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return BaseEvent{}, fmt.Errorf("event parsing error: %w", err)
	}

	return event, nil
}

// ParsePayload event payload'ını belirtilen yapıya dönüştürür
func (e *BaseEvent) ParsePayload(target interface{}) error {
	return json.Unmarshal(e.Payload, target)
} 