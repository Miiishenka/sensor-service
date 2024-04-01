package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"sync"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventIsNil    = errors.New("event is nil")
)

type EventRepository struct {
	events sync.Map
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if event == nil {
			return ErrEventIsNil
		}
		event.Timestamp = time.Now()
		r.events.Store(event, struct{}{})
		return nil
	}
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var lastEvent *domain.Event
		var lastEventTime time.Time
		r.events.Range(func(key, _ any) bool {
			event, ok := key.(*domain.Event)
			if !ok || event.SensorID != id {
				return true
			}

			if event.Timestamp.After(lastEventTime) {
				lastEventTime = event.Timestamp
				lastEvent = event
			}

			return true
		})

		if lastEvent == nil {
			return nil, ErrEventNotFound
		}

		return lastEvent, nil
	}
}
