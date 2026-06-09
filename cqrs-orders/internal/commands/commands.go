package commands

// Commands are plain structs - no methods, no logic.
// They're just data carriers expressing intent

type PlaceOrder struct {
	OrderID    string
	CustomerID string
	Product    string
	Quantity   int
}

type ConfirmOrder struct {
	OrderID string
}

type CancelOrder struct {
	OrderID string
}
