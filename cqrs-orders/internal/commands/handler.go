package commands

import (
	"cqrs-orders/internal/domain"
	"cqrs-orders/internal/eventstore"
	"time"
)

// Handler processes all write-side commands.
// It depends on the event store (to append events) and a read model
// (to check current state).
type OrderReader interface {
	GetByID(id string) (*domain.OrderState, bool)
}

// Add a Publisher to the Handler struct
type Publisher interface {
	Publish(e domain.Event)
}

type Handler struct {
	store     *eventstore.Store
	orders    OrderReader
	publisher Publisher
}

func NewHandler(store *eventstore.Store, orders OrderReader, publisher Publisher) *Handler {
	return &Handler{store: store, orders: orders, publisher: publisher}
}

func (h *Handler) HandlePlaceOrder(cmd PlaceOrder) error {
	// input validation
	if cmd.OrderID == "" || cmd.CustomerID == "" || cmd.Product == "" {
		return domain.ErrInvalidInput
	}
	if cmd.Quantity <= 0 {
		return domain.ErrInvalidQuantity
	}
	// model validation - check invariants
	// emit event
	h.emit(domain.Event{
		Type:       domain.EventOrderPlaced,
		OccurredAt: time.Now(),
		Payload: domain.OrderPlacedPayload{
			OrderID:    cmd.OrderID,
			CustomerID: cmd.CustomerID,
			Product:    cmd.Product,
			Quantity:   cmd.Quantity,
		},
	})
	return nil
}

func (h *Handler) HandleConfirmOrder(cmd ConfirmOrder) error {
	// load current state
	order, ok := h.orders.GetByID(cmd.OrderID)
	if !ok {
		return domain.ErrNotFound
	}
	// enforce business rules
	if order.Status == domain.StatusConfirmed {
		return domain.ErrAlreadyConfirmed
	}
	if order.Status == domain.StatusCancelled {
		return domain.ErrAlreadyCancelled
	}
	// emit event
	h.emit(domain.Event{
		Type:       domain.EventOrderConfirmed,
		OccurredAt: time.Now(),
		Payload: domain.OrderConfirmedPayload{
			OrderID: cmd.OrderID,
		},
	})
	return nil
}

func (h *Handler) HandleCancelOrder(cmd CancelOrder) error {
	order, ok := h.orders.GetByID(cmd.OrderID)
	if !ok {
		return domain.ErrNotFound
	}
	if order.Status == domain.StatusCancelled {
		return domain.ErrAlreadyCancelled
	}
	h.emit(domain.Event{
		Type:       domain.EventOrderCancelled,
		OccurredAt: time.Now(),
		Payload: domain.OrderCancelledPayload{
			OrderID: cmd.OrderID,
		},
	})
	return nil
}

// emit appends to the event store.
// In a real system, this is also where you'd publish to a message bus
// (Kafka, NATS, RabitMQ) so async consumers can update their own projections
func (h *Handler) emit(e domain.Event) {
	h.store.Append(e)
	h.publisher.Publish(e)
}
