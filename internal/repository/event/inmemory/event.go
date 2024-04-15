package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

var ErrEventIsNil = errors.New("event is nil")

type EventRepository struct {
	events sync.Map
}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (r *EventRepository) getEvents(sensorId int64) []*domain.Event {
	owners, ok := r.events.Load(sensorId)
	if !ok {
		return []*domain.Event{}
	}

	return owners.([]*domain.Event)
}

var mu sync.Mutex

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if event == nil {
			return ErrEventIsNil
		}

		mu.Lock()
		r.events.Store(event.SensorID, append(r.getEvents(event.SensorID), event))
		mu.Unlock()
		return nil
	}
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, ok := r.events.Load(id)
		if !ok {
			return nil, usecase.ErrEventNotFound
		}

		events, ok := value.([]*domain.Event)
		if !ok || len(events) < 1 {
			return nil, usecase.ErrEventNotFound
		}

		lastEvent := events[0]
		for _, event := range events {
			if event.Timestamp.After(lastEvent.Timestamp) {
				lastEvent = event
			}
		}

		return lastEvent, nil
	}
}
