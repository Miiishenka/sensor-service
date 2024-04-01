package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	ur  UserRepository
	sor SensorOwnerRepository
	sr  SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{ur: ur, sor: sor, sr: sr}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.Name == "" {
		return nil, ErrInvalidUserName
	}

	return user, u.ur.SaveUser(ctx, user)
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	_, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sr.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.sor.SaveSensorOwner(ctx, domain.SensorOwner{UserID: userID, SensorID: sensorID})
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	_, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	owners, err := u.sor.GetSensorsByUserID(ctx, userID)
	sensors := make([]domain.Sensor, 0, len(owners))

	for _, owner := range owners {
		sensor, findErr := u.sr.GetSensorByID(ctx, owner.SensorID)
		if findErr != nil {
			return nil, findErr
		}
		sensors = append(sensors, *sensor)
	}

	return sensors, err
}
