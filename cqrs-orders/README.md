# CQRS 

(Command Query Responsibility Segregation) is a design pattern that separates the read and write operations of a system into distinct models. This allows for better scalability, maintainability, and performance by optimizing each side for its specific responsibilities.

### Structure

```
domain/     → entities, errors, events
eventstore/ → append-only log
eventbus/   → async pub/sub between write and read sides
readmodel/  → projection, updated via bus subscription
commands/   → publishes to bus
queries/    → pure reads
```


### Test Output

```bash
dev/CloudDesignPatterns/cqrs-orders via 🐹 v1.26.4-X:nodwarf5 
❯ go run ./...
=== WRITE SIDE: Commands ===
  ✓ PlaceOrder(order-1, Laptop)
  ✓ PlaceOrder(order-2, Mouse)
  ✓ PlaceOrder(order-1, Keyboard)
  ✓ ConfirmOrder(order-1)
  ✓ CancelOrder(order-2)
  ✗ ConfirmOrder(order-2) [cancelled - should fail] → order already cancelled
  ✗ ConfirmOrder(order-1) [already confirmed - should fail] → order already confirmed
  ✗ PlaceOrder(bad quantity) [should fail] → invalid quantity

=== READ SIDE: Queries ===
  ✓ GetOrder(order-1) → Laptop x1 (confirmed)

  ✓ ListOrders() → 3 orders
    - order-1: Laptop x1 (confirmed)
    - order-2: Mouse x3 (cancelled)
    - order-3: Keyboard x2 (pending)

  ✓ ListOrders(status=pending) → 1 orders
    - order-3: Keyboard x2 (pending)

  ✗ GetOrder(ghost) → order not found
    errors.Is(ErrNotFound) → true

=== EVENT STORE: Audit Log ===
   #1 [OrderPlaced] at 10:15:54.179
   #2 [OrderPlaced] at 10:15:54.179
   #3 [OrderPlaced] at 10:15:54.179
   #4 [OrderConfirmed] at 10:15:54.179
   #5 [OrderCancelled] at 10:15:54.179
```
