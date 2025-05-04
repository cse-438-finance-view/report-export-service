package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/burakmike/report-export-service/pkg/event"
)

// stubHandler is a test double for EventHandler
// It records if Handle was called and with which event, and can return a predetermined error.
type stubHandler struct {
	Et      event.EventType
	Handled bool
	Recv    event.BaseEvent
	RetErr  error
}

func (s *stubHandler) EventType() event.EventType {
	return s.Et
}

func (s *stubHandler) Handle(ctx context.Context, evt event.BaseEvent) error {
	s.Handled = true
	s.Recv = evt
	return s.RetErr
}

func TestRegisterAndGetHandler(t *testing.T) {
	registry := NewHandlerRegistry()
	h1 := &stubHandler{Et: event.PortfolioReport}
	registry.RegisterHandler(h1)

	got := registry.GetHandler(event.PortfolioReport)
	if got != h1 {
		t.Errorf("GetHandler returned %v, want %v", got, h1)
	}

	if registry.GetHandler("unknown.event") != nil {
		t.Errorf("GetHandler for unknown event should return nil")
	}
}

func TestHandleEvent_Route(t *testing.T) {
	registry := NewHandlerRegistry()
	h1 := &stubHandler{Et: event.PortfolioReport}
	h2 := &stubHandler{Et: "other.event"}
	registry.RegisterHandler(h1)
	registry.RegisterHandler(h2)

	evt := event.BaseEvent{EventType: event.PortfolioReport, Timestamp: "x", Payload: nil}
	err := registry.HandleEvent(context.Background(), evt)
	if err != nil {
		t.Fatalf("HandleEvent returned unexpected error: %v", err)
	}
	if !h1.Handled {
		t.Error("Expected h1 to handle event")
	}
	if h2.Handled {
		t.Error("Did not expect h2 to handle event")
	}
}

func TestHandleEvent_Error(t *testing.T) {
	registry := NewHandlerRegistry()
	sentinel := errors.New("handler error")
	h := &stubHandler{Et: event.PortfolioReport, RetErr: sentinel}
	registry.RegisterHandler(h)

	evt := event.BaseEvent{EventType: event.PortfolioReport, Timestamp: "x", Payload: nil}
	err := registry.HandleEvent(context.Background(), evt)
	if err != sentinel {
		t.Errorf("HandleEvent error = %v, want %v", err, sentinel)
	}
} 