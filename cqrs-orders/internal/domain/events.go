package domain

import "time"

// EventType is a named string so the compiler catches typos
type EventType string

const (
	EventOrderPlaced    EventType = "OrderPlaced"
	EventOrderConfirmed EventType = "OrderConfirmed"
	EventOrderCancelled EventType = "OrderCancelled"
)

// Event is an envelope - every event gets wrapped in this.
// The payload holds the event-specific data
type Event struct {
	Type       EventType
	OccurredAt time.Time
	Payload    any
}

// --- Payloads ---
// Each payload carries exactly the data that changed.
// Think of these as the "what happened" record

type OrderPlacedPayload struct {
	OrderID    string
	CustomerID string
	Product    string
	Quantity   int
}

type OrderConfirmedPayload struct {
	OrderID string
}

type OrderCancelledPayload struct {
	OrderID string
}
