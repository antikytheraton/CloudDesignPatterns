package domain

import "errors"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
)

// Order is our aggregate - the central business object.
// Notice it has No DB tags, No JSON tags, No infra concerns
// It's pure business logic
type Order struct {
	ID         string
	CustomerID string
	Product    string
	Quantity   int
	Status     OrderStatus
}

// OrderState is the minimal state the handler needs to enforce rules.
// it's not the read model view - it's just enough to answer invariant questions
type OrderState struct {
	Status OrderStatus
}

// Domain errors live here too - they're part of the business language
var (
	ErrNotFound         = errors.New("order not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrAlreadyCancelled = errors.New("order already cancelled")
	ErrAlreadyConfirmed = errors.New("order already confirmed")
)
