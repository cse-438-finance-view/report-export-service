package handler

import (
	"context"

	"github.com/burakmike/report-export-service/pkg/event"
)

// EventHandler bir olayı işleyen arayüz
type EventHandler interface {
	// EventType bu işleyicinin hangi event tipini işlediğini belirtir
	EventType() event.EventType
	
	// Handle olayı işler
	Handle(ctx context.Context, evt event.BaseEvent) error
}

// HandlerRegistry farklı event tipleri için farklı işleyicileri yönetir
type HandlerRegistry struct {
	handlers map[event.EventType]EventHandler
}

// NewHandlerRegistry yeni bir işleyici kaydı oluşturur
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[event.EventType]EventHandler),
	}
}

// RegisterHandler belirli bir event tipi için bir işleyici ekler
func (r *HandlerRegistry) RegisterHandler(handler EventHandler) {
	r.handlers[handler.EventType()] = handler
}

// GetHandler belirli bir event tipi için işleyici döndürür, bulunamazsa nil döner
func (r *HandlerRegistry) GetHandler(eventType event.EventType) EventHandler {
	handler, exists := r.handlers[eventType]
	if !exists {
		return nil
	}
	return handler
}

// HandleEvent bir olayı uygun işleyiciye yönlendirir
func (r *HandlerRegistry) HandleEvent(ctx context.Context, evt event.BaseEvent) error {
	handler := r.GetHandler(evt.EventType)
	if handler == nil {
		return nil // İlgili olay tipi için işleyici bulunamadı
	}
	
	return handler.Handle(ctx, evt)
} 