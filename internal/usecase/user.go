package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	userRepository   UserRepository
	ownerRepository  SensorOwnerRepository
	sensorRepository SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{userRepository: ur, ownerRepository: sor, sensorRepository: sr}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.Name == "" {
		return nil, ErrInvalidUserName
	}

	return user, u.userRepository.SaveUser(ctx, user)
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	_, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sensorRepository.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.ownerRepository.SaveSensorOwner(ctx, domain.SensorOwner{UserID: userID, SensorID: sensorID})
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	_, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	owners, err := u.ownerRepository.GetSensorsByUserID(ctx, userID)
	var sensors []domain.Sensor

	for _, owner := range owners {
		sensor, findErr := u.sensorRepository.GetSensorByID(ctx, owner.SensorID)
		if findErr != nil {
			return nil, findErr
		}
		sensors = append(sensors, *sensor)
	}

	return sensors, err
}
