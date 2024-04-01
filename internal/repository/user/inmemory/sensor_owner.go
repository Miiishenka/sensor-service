package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	sensorOwners sync.Map
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		owners, err := r.GetSensorsByUserID(ctx, sensorOwner.UserID)
		if err != nil {
			return err
		}

		for _, owner := range owners {
			if owner.SensorID == sensorOwner.SensorID {
				return nil
			}
		}

		r.sensorOwners.Store(sensorOwner.UserID, append(owners, sensorOwner))
		return nil
	}
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		owners, ok := r.sensorOwners.Load(userID)
		if !ok {
			return []domain.SensorOwner{}, nil
		}

		return owners.([]domain.SensorOwner), nil
	}
}
