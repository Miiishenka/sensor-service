package inmemory

import (
	"context"
	"errors"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"sync/atomic"
)

var ErrUserIsNull = errors.New("user is nil")

type UserRepository struct {
	users  sync.Map
	lastId int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if user == nil {
			return ErrUserIsNull
		}
		user.ID = atomic.AddInt64(&r.lastId, 1)
		r.users.Store(user.ID, user)
		return nil
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		user, ok := r.users.Load(id)
		if !ok {
			return nil, usecase.ErrUserNotFound
		}

		return user.(*domain.User), nil
	}
}
