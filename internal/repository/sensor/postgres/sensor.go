package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/domain"
	"homework/internal/usecase"
	"time"
)

const (
	saveSensorQuery = `
		INSERT INTO sensors 
    	(serial_number, type, current_state, description, is_active, registered_at, last_activity) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	updateSensorQuery = `
		UPDATE sensors
		SET type=$3, current_state=$4, description=$5, is_active=$6, last_activity=$7
    	WHERE id=$1 AND serial_number=$2`

	getSensorsQuery              = `SELECT * FROM sensors`
	getSensorByIDQuery           = `SELECT * FROM sensors WHERE id=$1`
	getSensorBySerialNumberQuery = `SELECT * FROM sensors WHERE serial_number=$1`
)

type SensorRepository struct {
	pool *pgxpool.Pool
}

func NewSensorRepository(pool *pgxpool.Pool) *SensorRepository {
	return &SensorRepository{
		pool: pool,
	}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if sensor.ID == 0 {
		sensor.RegisteredAt = time.Now()
		return r.pool.QueryRow(ctx, saveSensorQuery, sensor.SerialNumber, sensor.Type, sensor.CurrentState,
			sensor.Description, sensor.IsActive, sensor.RegisteredAt, sensor.LastActivity).Scan(&sensor.ID)
	}

	_, err := r.pool.Exec(ctx, updateSensorQuery, sensor.ID, sensor.SerialNumber, sensor.Type, sensor.CurrentState,
		sensor.Description, sensor.IsActive, sensor.LastActivity)

	return err
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	rows, err := r.pool.Query(ctx, getSensorsQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		sensor := domain.Sensor{}
		if err := rows.Scan(&sensor.ID, &sensor.SerialNumber, &sensor.Type, &sensor.CurrentState, &sensor.Description,
			&sensor.IsActive, &sensor.RegisteredAt, &sensor.LastActivity); err != nil {
			return nil, err
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	return r.getSensorBy(ctx, getSensorByIDQuery, id)
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	return r.getSensorBy(ctx, getSensorBySerialNumberQuery, sn)
}

func (r *SensorRepository) getSensorBy(ctx context.Context, query string, key any) (*domain.Sensor, error) {
	row := r.pool.QueryRow(ctx, query, key)
	sensor := &domain.Sensor{}
	err := row.Scan(&sensor.ID, &sensor.SerialNumber, &sensor.Type, &sensor.CurrentState, &sensor.Description,
		&sensor.IsActive, &sensor.RegisteredAt, &sensor.LastActivity)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrSensorNotFound
	}

	if err != nil {
		return nil, err
	}

	return sensor, nil
}
