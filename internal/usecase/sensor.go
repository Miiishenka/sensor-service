package usecase

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/repository/sensor/inmemory"
)

type Sensor struct {
	sensorRepository SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{sensorRepository: sr}
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	if sensor.Type != domain.SensorTypeADC && sensor.Type != domain.SensorTypeContactClosure {
		return nil, ErrWrongSensorType
	}

	if len(sensor.SerialNumber) != 10 {
		return nil, ErrWrongSensorSerialNumber
	}

	findSensor, err := s.sensorRepository.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
	if errors.Is(err, inmemory.ErrSensorNotFound) {
		return sensor, s.sensorRepository.SaveSensor(ctx, sensor)
	}

	if err != nil {
		return nil, err
	}

	return findSensor, err
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	return s.sensorRepository.GetSensors(ctx)
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	return s.sensorRepository.GetSensorByID(ctx, id)
}
