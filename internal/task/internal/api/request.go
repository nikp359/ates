package api

type (
	TasksShuffledResponse struct {
		Shuffled int `json:"shuffled"`
	}

	TaskCompletedRequest struct {
		TaskID string `json:"task_id"`
	}
)
