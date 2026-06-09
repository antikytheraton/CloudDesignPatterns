package queries

import "cqrs-orders/internal/domain"

// Queries are plain structs, just like commands.
// Names as nouns or questions, never imperatives.

type GetOrder struct {
	OrderID string
}

type ListOrders struct {
	Status domain.OrderStatus // optional filter - empty means all
}
