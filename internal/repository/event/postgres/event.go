package postgres

import (
	"context"
	"errors"
	"homework/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	saveEventQuery    = `INSERT INTO events (timestamp, sensor_serial_number, sensor_id, payload) VALUES ($1, $2, $3, $4)`
	getLastEventQuery = `SELECT * FROM events WHERE sensor_id=$1 ORDER BY timestamp DESC LIMIT 1`
	getSensorHistory  = `SELECT * FROM events WHERE sensor_id=$1 AND timestamp BETWEEN $2 AND $3`
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool,
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	_, err := r.pool.Exec(ctx, saveEventQuery, event.Timestamp, event.SensorSerialNumber, event.SensorID, event.Payload)
	return err
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	row := r.pool.QueryRow(ctx, getLastEventQuery, id)

	event := &domain.Event{}
	err := row.Scan(&event.Timestamp, &event.SensorSerialNumber, &event.SensorID, &event.Payload)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEventNotFound
	}

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) GetSensorHistory(ctx context.Context, sensorId int64, start, end time.Time) ([]domain.Event, error) {
	rows, err := r.pool.Query(ctx, getSensorHistory, sensorId, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event := domain.Event{}
		if err := rows.Scan(&event.Timestamp, &event.SensorSerialNumber, &event.SensorID, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
