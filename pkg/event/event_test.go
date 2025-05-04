package event

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewBaseEvent(t *testing.T) {
	payload := map[string]string{"foo": "bar"}
	evt, err := NewBaseEvent(PortfolioReport, payload)
	if err != nil {
		t.Fatalf("NewBaseEvent error: %v", err)
	}
	if evt.EventType != PortfolioReport {
		t.Errorf("EventType = %v; want %v", evt.EventType, PortfolioReport)
	}
	// Timestamp must be RFC3339
	if _, err := time.Parse(time.RFC3339, evt.Timestamp); err != nil {
		t.Errorf("Timestamp not RFC3339: %q", evt.Timestamp)
	}
	// Payload contains correct JSON
	var m map[string]string
	if err := json.Unmarshal(evt.Payload, &m); err != nil {
		t.Errorf("Unmarshal payload error: %v", err)
	} else if m["foo"] != "bar" {
		t.Errorf("payload foo = %q; want %q", m["foo"], "bar")
	}
}

func TestParseEvent_ValidAndInvalid(t *testing.T) {
	valid := []byte(`{"event_type":"portfolio.report","timestamp":"2006-01-02T15:04:05Z","payload":{"portfolios":[]}}`)
	evt, err := ParseEvent(valid)
	if err != nil {
		t.Fatalf("ParseEvent valid error: %v", err)
	}
	if evt.EventType != PortfolioReport {
		t.Errorf("EventType = %v; want %v", evt.EventType, PortfolioReport)
	}
	// invalid JSON should error
	if _, err := ParseEvent([]byte(`invalid`)); err == nil {
		t.Errorf("ParseEvent invalid did not return error")
	}
}

func TestParsePayload(t *testing.T) {
	raw := json.RawMessage(`{"name":"test"}`)
	evt := BaseEvent{EventType: PortfolioReport, Timestamp: "x", Payload: raw}
	var m map[string]string
	if err := evt.ParsePayload(&m); err != nil {
		t.Fatalf("ParsePayload error: %v", err)
	}
	if m["name"] != "test" {
		t.Errorf("Parsed name = %q; want %q", m["name"], "test")
	}
} 