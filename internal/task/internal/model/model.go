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

	Task struct {
		PublicID       string    `json:"public_id" db:"public_id"`
		Title          string    `json:"title" db:"title"`
		JiraID         string    `json:"jira_id" db:"jira_id"`
		Description    string    `json:"description" db:"description"`
		Status         string    `json:"status" db:"status"`
		AssignedUserID string    `json:"assigned_user_id" db:"assigned_user_id"`
		UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	}
)
