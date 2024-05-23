package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/domain"
	"homework/internal/usecase"
)

const (
	saveUserQuery    = `INSERT INTO users (name) VALUES ($1) RETURNING id`
	getUserByIDQuery = `SELECT * FROM users WHERE id=$1`
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	return r.pool.QueryRow(ctx, saveUserQuery, user.Name).Scan(&user.ID)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	err := r.pool.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
