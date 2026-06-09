package readmodel

import (
	"cqrs-orders/internal/domain"
	"sync"
	"time"
)

// OrderView is the query-optimized shape of an order.
// Notice it's flat - no nested structs, No IDs to join on.
// Everything a query needs is already here.
type OrderView struct {
	ID         string
	CustomerID string
	Product    string
	Quantity   int
	Status     domain.OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// OrderStore is the in-memory read-side store.
// It satisfices the commands.OrderReader interface automatically
// because it implements GetByID - Go's implicit interface satisfaction.
type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]*OrderView
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]*OrderView),
	}
}

// --- commands.OrderReader interface ---

// GetByID returns the minimal state the command handler needs.
// It returns a value type, not a pointer, so the handler
// can't accidentally mutate the read model though it.
func (s *OrderStore) GetByID(id string) (*domain.OrderState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, false
	}
	return &domain.OrderState{Status: o.Status}, true
}

func (s *OrderStore) Apply(e domain.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch e.Type {
	case domain.EventOrderPlaced:
		p := e.Payload.(domain.OrderPlacedPayload)
		s.orders[p.OrderID] = &OrderView{
			ID:         p.OrderID,
			CustomerID: p.CustomerID,
			Product:    p.Product,
			Quantity:   p.Quantity,
			Status:     domain.StatusPending,
			CreatedAt:  e.OccurredAt,
			UpdatedAt:  e.OccurredAt,
		}

	case domain.EventOrderConfirmed:
		p := e.Payload.(domain.OrderConfirmedPayload)
		if o, ok := s.orders[p.OrderID]; ok {
			o.Status = domain.StatusConfirmed
			o.UpdatedAt = e.OccurredAt
		}

	case domain.EventOrderCancelled:
		p := e.Payload.(domain.OrderCancelledPayload)
		if o, ok := s.orders[p.OrderID]; ok {
			o.Status = domain.StatusCancelled
			o.UpdatedAt = e.OccurredAt
		}
	}
}

// --- Query methods ---

func (s *OrderStore) FindByID(id string) (*OrderView, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	return o, ok
}

func (s *OrderStore) FindAll() []*OrderView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*OrderView, 0, len(s.orders))
	for _, o := range s.orders {
		result = append(result, o)
	}
	return result
}

func (s *OrderStore) FindByStatus(status domain.OrderStatus) []*OrderView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*OrderView
	for _, o := range s.orders {
		if o.Status == status {
			result = append(result, o)
		}
	}
	return result
}
