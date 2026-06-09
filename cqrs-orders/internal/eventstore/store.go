package eventstore

import (
	"cqrs-orders/internal/domain"
	"sync"
)

// Store is an in-memory, append-only event log.
// In production this would be PostgreSQL, EventStoreDB, Kafka, etc.
// The interface stays identical - only the backing storage changes
type Store struct {
	mu     sync.RWMutex
	events []domain.Event
}

func New() *Store {
	return &Store{
		events: make([]domain.Event, 0),
	}
}

func (s *Store) Append(e domain.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
}

func (s *Store) GetAll() []domain.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Event, len(s.events))
	copy(result, s.events)
	return result
}
