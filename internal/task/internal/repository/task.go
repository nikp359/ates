package repository

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
	"github.com/nikp359/ates/internal/task/internal/model"
)

const (
	insertTask = `
INSERT IGNORE INTO task (public_id, title, jira_id, description, assigned_user_id)
    VALUES (:public_id, :title, :jira_id, :description,
            (SELECT public_id FROM user AS usr1 JOIN
                (SELECT (RAND() *
                    (SELECT MAX(id) FROM user WHERE role='employee')) AS user_id) AS usr2
                        WHERE usr1.id >= usr2.user_id AND usr1.role='employee' ORDER BY usr1.id LIMIT 1)
            )
`
	selectTask = `SELECT public_id, title, jira_id, description, status, assigned_user_id, updated_at FROM task
WHERE public_id=?;`

	selectNewTasks = `SELECT public_id, title, jira_id, description, status, assigned_user_id, updated_at FROM task WHERE status='new';`

	selectUsers = `SELECT public_id FROM user;`

	updateTaskAssignedUser = `UPDATE task SET assigned_user_id=:assigned_user_id WHERE public_id=:public_id;`

	updateTaskStatus = `UPDATE task SET status=? WHERE public_id=?;`

	TaskStatusNew       = "new"
	TaskStatusCompleted = "completed"
)

type (
	TaskRepository struct {
		db *sqlx.DB
	}
)

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) GetByPublicID(publicID string) (model.Task, error) {
	var task model.Task
	if err := r.db.Get(&task, selectTask, publicID); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Add(task *model.Task) error {
	if _, err := r.db.NamedExec(insertTask, task); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Shuffled(ctx context.Context) ([]model.Task, error) {
	// TODO: move to transactions to prevent a data race
	userIDs := make([]string, 0)
	if err := r.db.SelectContext(ctx, &userIDs, selectUsers); err != nil {
		return nil, err
	}

	tasks := make([]model.Task, 0)
	if err := r.db.SelectContext(ctx, &tasks, selectNewTasks); err != nil {
		return nil, err
	}

	for i := range tasks {
		tasks[i].AssignedUserID = gerRandomUserID(userIDs)

		_, err := r.db.NamedExecContext(ctx, updateTaskAssignedUser, tasks[i])
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func (r *TaskRepository) ChangeStatus(ctx context.Context, taskID string, status string) (model.Task, error) {
	if _, err := r.db.ExecContext(ctx, updateTaskStatus, status, taskID); err != nil {
		return model.Task{}, fmt.Errorf("update status taskID: %s err: %w", taskID, err)
	}

	task, err := r.GetByPublicID(taskID)
	if err != nil {
		return model.Task{}, fmt.Errorf("get task, taskID: %s err: %w", taskID, err)
	}

	return task, nil
}

func gerRandomUserID(userIDs []string) string {
	return userIDs[rand.Intn(len(userIDs))]
}
