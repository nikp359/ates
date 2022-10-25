package model

import "time"

type (
	User struct {
		PublicID  string    `db:"public_id"`
		Email     string    `db:"email"`
		Role      string    `db:"role"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)
