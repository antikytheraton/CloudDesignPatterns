package queries

import (
	"cqrs-orders/internal/domain"
	"cqrs-orders/internal/readmodel"
)

type Store interface {
	FindByID(id string) (*readmodel.OrderView, bool)
	FindAll() []*readmodel.OrderView
	FindByStatus(status domain.OrderStatus) []*readmodel.OrderView
}

type Handler struct {
	orders Store
}

func NewHandler(orders Store) *Handler {
	return &Handler{orders: orders}
}

func (h *Handler) HandleGetOrder(q GetOrder) (*readmodel.OrderView, error) {
	o, ok := h.orders.FindByID(q.OrderID)
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (h *Handler) HandleListOrders(q ListOrders) ([]*readmodel.OrderView, error) {
	if q.Status != "" {
		return h.orders.FindByStatus(q.Status), nil
	}
	return h.orders.FindAll(), nil
}
