package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/nikp359/ates/internal/auth/internal/model"
)

type (
	UserRepository struct {
		db *sqlx.DB
	}
)

const insertUser = `INSERT INTO user (public_id, email, role) VALUES (:public_id, :email, :role);`

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(user *model.User) error {
	if _, err := r.db.NamedExec(insertUser, user); err != nil {
		return err
	}

	return nil
}
