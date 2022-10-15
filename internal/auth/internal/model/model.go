package model

type (
	User struct {
		PublicID string `json:"public_id" db:"public_id"`
		Email    string `json:"email" db:"email"`
		Role     string `json:"role" db:"role"`
	}
)
