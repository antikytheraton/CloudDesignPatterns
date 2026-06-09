package main

import (
	"cqrs-orders/internal/commands"
	"cqrs-orders/internal/domain"
	"cqrs-orders/internal/eventbus"
	"cqrs-orders/internal/eventstore"
	"cqrs-orders/internal/queries"
	"cqrs-orders/internal/readmodel"
	"errors"
	"fmt"
)

func main() {
	// --- composition root ---
	store := eventstore.New()
	orderStore := readmodel.NewOrderStore()
	bus := eventbus.New(100)

	bus.Subscribe(orderStore.Apply)

	bus.Start()
	defer bus.Stop()

	cmdHandler := commands.NewHandler(store, orderStore, bus)
	qryHandler := queries.NewHandler(orderStore)

	fmt.Printf("=== WRITE SIDE: Commands ===\n")

	run("PlaceOrder(order-1, Laptop)",
		cmdHandler.HandlePlaceOrder(commands.PlaceOrder{
			OrderID:    "order-1",
			CustomerID: "customer-A",
			Product:    "Laptop",
			Quantity:   1,
		}))
	run("PlaceOrder(order-2, Mouse)",
		cmdHandler.HandlePlaceOrder(commands.PlaceOrder{
			OrderID:    "order-2",
			CustomerID: "customer-B",
			Product:    "Mouse",
			Quantity:   3,
		}))
	run("PlaceOrder(order-1, Keyboard)",
		cmdHandler.HandlePlaceOrder(commands.PlaceOrder{
			OrderID:    "order-3",
			CustomerID: "customer-A",
			Product:    "Keyboard",
			Quantity:   2,
		}))
	run("ConfirmOrder(order-1)",
		cmdHandler.HandleConfirmOrder(commands.ConfirmOrder{OrderID: "order-1"}))
	run("CancelOrder(order-2)",
		cmdHandler.HandleCancelOrder(commands.CancelOrder{OrderID: "order-2"}))

	// Business rule violations
	run("ConfirmOrder(order-2) [cancelled - should fail]",
		cmdHandler.HandleConfirmOrder(commands.ConfirmOrder{OrderID: "order-2"}))
	run("ConfirmOrder(order-1) [already confirmed - should fail]",
		cmdHandler.HandleConfirmOrder(commands.ConfirmOrder{OrderID: "order-1"}))
	run("PlaceOrder(bad quantity) [should fail]",
		cmdHandler.HandlePlaceOrder(commands.PlaceOrder{
			OrderID:    "order-4",
			CustomerID: "customer-C",
			Product:    "Chair",
			Quantity:   -1,
		}))

	fmt.Printf("\n=== READ SIDE: Queries ===\n")

	// Single order
	order, err := qryHandler.HandleGetOrder(queries.GetOrder{OrderID: "order-1"})
	if err != nil {
		fmt.Printf("  ✗ GetOrder(order-1) → %v\n", err)
	} else {
		fmt.Printf("  ✓ GetOrder(order-1) → %s x%d (%s)\n", order.Product, order.Quantity, order.Status)
	}

	// All orders
	all, _ := qryHandler.HandleListOrders(queries.ListOrders{})
	fmt.Printf("\n  ✓ ListOrders() → %d orders\n", len(all))
	for _, o := range all {
		fmt.Printf("    - %s: %s x%d (%s)\n", o.ID, o.Product, o.Quantity, o.Status)
	}

	// Filtered
	pending, _ := qryHandler.HandleListOrders(queries.ListOrders{Status: "pending"})
	fmt.Printf("\n  ✓ ListOrders(status=pending) → %d orders\n", len(pending))
	for _, o := range pending {
		fmt.Printf("    - %s: %s x%d (%s)\n", o.ID, o.Product, o.Quantity, o.Status)
	}

	// Not found
	_, err = qryHandler.HandleGetOrder(queries.GetOrder{OrderID: "ghost"})
	fmt.Printf("\n  ✗ GetOrder(ghost) → %v\n", err)
	fmt.Printf("    errors.Is(ErrNotFound) → %v\n", errors.Is(err, domain.ErrNotFound))

	fmt.Printf("\n=== EVENT STORE: Audit Log ===\n")
	for i, e := range store.GetAll() {
		fmt.Printf("   #%d [%s] at %s\n", i+1, e.Type, e.OccurredAt.Format("15:04:05.000"))
	}
}

func run(label string, err error) {
	if err != nil {
		fmt.Printf("  ✗ %s → %v\n", label, err)
	} else {
		fmt.Printf("  ✓ %s\n", label)
	}
}
