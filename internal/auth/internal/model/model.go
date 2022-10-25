package model

import "time"

type (
	User struct {
		PublicID  string    `json:"public_id" db:"public_id"`
		Email     string    `json:"email" db:"email"`
		Role      string    `json:"role" db:"role"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}
)
