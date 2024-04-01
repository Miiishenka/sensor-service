package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrSensorNotFound = errors.New("sensor not found")
	ErrSensorIsNil    = errors.New("sensor is nil")
)

type SensorRepository struct {
	sensors sync.Map
	lastId  int64
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if sensor == nil {
			return ErrSensorIsNil
		}

		existSensor, _ := r.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
		if existSensor != nil {
			return nil
		}

		sensor.ID = atomic.AddInt64(&r.lastId, 1)
		sensor.RegisteredAt = time.Now()
		r.sensors.Store(sensor.ID, sensor)
		return nil
	}
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var sensors []domain.Sensor
		r.sensors.Range(func(_, value any) bool {
			sensors = append(sensors, *value.(*domain.Sensor))
			return true
		})

		return sensors, nil
	}
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		sensor, ok := r.sensors.Load(id)
		if !ok {
			return nil, ErrSensorNotFound
		}

		return sensor.(*domain.Sensor), nil
	}
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var foundSensor *domain.Sensor
		var found bool
		r.sensors.Range(func(_, value any) bool {
			sensor, ok := value.(*domain.Sensor)
			if !ok || sensor.SerialNumber != sn {
				return true
			}
			foundSensor = sensor
			found = true
			return false
		})

		if !found {
			return nil, ErrSensorNotFound
		}

		return foundSensor, nil
	}
}
