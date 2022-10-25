package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/nikp359/ates/internal/task/internal/model"
)

type (
	UserRepository struct {
		db *sqlx.DB
	}
)

const (
	insertUser = `INSERT IGNORE INTO user (public_id, email, role, updated_at) VALUES (:public_id, :email, :role, :updated_at);`

	selectUser = `SELECT public_id, email, role, created_at, updated_at FROM user WHERE public_id=?;`

	updateUser = `
	INSERT IGNORE INTO user (public_id, email, role, updated_at)
    	VALUES (:public_id, :email, :role, :updated_at)
		ON DUPLICATE KEY UPDATE role=:role, updated_at=:updated_at;`

	deleteUser = `DELETE FROM user where public_id=?;`
)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByPublicID(publicID string) (*model.User, error) {
	var user model.User
	if err := r.db.Get(&user, selectUser, publicID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Add(user *model.User) error {
	if _, err := r.db.NamedExec(insertUser, user); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *model.User) error {
	if _, err := r.db.NamedExec(updateUser, user); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(publicID string) error {
	if _, err := r.db.Exec(deleteUser, publicID); err != nil {
		return err
	}

	return nil
}
