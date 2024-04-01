package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/repository/sensor/inmemory"
)

type Sensor struct {
	sr SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{sr: sr}
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	if sensor.Type != domain.SensorTypeADC && sensor.Type != domain.SensorTypeContactClosure {
		return nil, ErrWrongSensorType
	}

	if len(sensor.SerialNumber) != 10 {
		return nil, ErrWrongSensorSerialNumber
	}

	findSensor, err := s.sr.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
	if errors.Is(err, inmemory.ErrSensorNotFound) {
		return sensor, s.sr.SaveSensor(ctx, sensor)
	}

	if err != nil {
		return nil, err
	}

	return findSensor, err
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	return s.sr.GetSensors(ctx)
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	return s.sr.GetSensorByID(ctx, id)
}
