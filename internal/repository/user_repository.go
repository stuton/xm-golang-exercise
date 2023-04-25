package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/stuton/xm-golang-exercise/internal/model"
)

//go:generate mockery --name UserRepository --filename user_repository.go
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepository{
		db: db,
	}
}

// GetOne implements Company
func (repo userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	u := new(model.User)

	if err := repo.db.Get(u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return nil, err
	}

	return u, nil
}
