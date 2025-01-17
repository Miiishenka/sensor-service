package usecase

import (
	"context"
	"homework/internal/domain"
	"time"
)

type Event struct {
	er EventRepository
	sr SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{er: er, sr: sr}
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if event.Timestamp.IsZero() {
		return ErrInvalidEventTimestamp
	}

	sensor, err := e.sr.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}

	event.SensorID = sensor.ID
	err = e.er.SaveEvent(ctx, event)
	if err != nil {
		return err
	}

	sensor.CurrentState = event.Payload
	sensor.LastActivity = event.Timestamp

	err = e.sr.SaveSensor(ctx, sensor)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	return e.er.GetLastEventBySensorID(ctx, id)
}

func (e *Event) GetSensorHistory(ctx context.Context, sensorId int64, start, end time.Time) ([]domain.Event, error) {
	return e.er.GetSensorHistory(ctx, sensorId, start, end)
}
