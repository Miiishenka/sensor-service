package postgres

import (
	"context"
	"homework/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	saveSensorOwnerQuery    = `INSERT INTO sensors_users (sensor_id, user_id) VALUES ($1, $2)`
	getSensorsByUserIdQuery = `SELECT sensor_id, user_id FROM sensors_users WHERE user_id=$1`
)

type SensorOwnerRepository struct {
	pool *pgxpool.Pool
}

func NewSensorOwnerRepository(pool *pgxpool.Pool) *SensorOwnerRepository {
	return &SensorOwnerRepository{
		pool,
	}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	_, err := r.pool.Exec(ctx, saveSensorOwnerQuery, sensorOwner.SensorID, sensorOwner.UserID)
	return err
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	rows, err := r.pool.Query(ctx, getSensorsByUserIdQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var owners []domain.SensorOwner
	for rows.Next() {
		owner := domain.SensorOwner{}
		if err := rows.Scan(&owner.SensorID, &owner.UserID); err != nil {
			return nil, err
		}
		owners = append(owners, owner)
	}

	return owners, nil
}
