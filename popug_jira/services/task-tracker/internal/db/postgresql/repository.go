package postgresql

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
--------------------------- Tasks ---------------------------
*/

func (i *Instance) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if err := i.database.WithContext(ctx).Create(task).Error; err != nil {
		log.Errorf(
			"i.database.WithContext(ctx).Create(%v): %s",
			task,
			err.Error(),
		)

		return nil, fmt.Errorf("error create task  %v in database: %w", task, err)
	}

	log.Debugf("success create task %v in database!", task)

	return task, nil
}

func (i *Instance) TaskCompleteByUser(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
) (*domain.Task, error) {
	task := &domain.Task{}

	err := i.database.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(
				"id = ? and user_id = ? and status = ?", taskID, userID, domain.TaskOpened,
			).First(task).Error; err != nil {
				return fmt.Errorf(
					"error find opened task %q assigned to user %q in db: %w",
					taskID.String(),
					userID.String(),
					err,
				)
			}

			task.Status = domain.TaskCompleted

			if err := tx.Updates(task).Error; err != nil {
				return fmt.Errorf(
					"error update task %q (%v) assigned to user %q in db: %w",
					taskID.String(),
					task,
					userID.String(),
					err,
				)
			}

			return nil
		})

	if err != nil {
		return nil, fmt.Errorf(
			"error complete task %q assigned to user %q in db: %w",
			taskID.String(),
			userID.String(),
			err,
		)
	}

	log.Debugf(
		"success complete task %q assigned to user %q in db!",
		taskID.String(),
		userID.String(),
	)

	return task, nil
}

func (i *Instance) GetAllTasksByUserAndStatus(
	ctx context.Context,
	userID uuid.UUID,
	status domain.TaskStatus,
) ([]*domain.Task, error) {
	var tasks []*domain.Task

	query := i.database.WithContext(ctx).Where("status = ?", status)
	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Find(&tasks).Error
	if err != nil {
		log.Errorf(
			"i.database.WithContext(ctx).Where(role = %s).Where(user_id = %s).Find(&users): %s",
			status,
			userID.String(),
			err.Error(),
		)

		return nil, fmt.Errorf("error get tasks by user id %q and status %q: %w", userID, status, err)
	}

	log.Debugf(
		"success get tasks by user id %q and status %q: %s", userID, status, err,
	)

	return tasks, nil
}

func (i *Instance) RandomlyUpdateAssignedOpenedTasks(
	ctx context.Context,
) (map[*domain.Task]*domain.User, error) {
	updatedTasks := make(map[*domain.Task]*domain.User)

	err := i.database.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			var workersDB = make([]domain.User, 0)

			if err := tx.Where("role = ?", domain.WorkerRole).
				Find(&workersDB).Error; err != nil {

				return fmt.Errorf("error get worker user: %w", err)
			}

			if len(workersDB) == 0 {
				return nil
			}

			var tasksDB = make([]domain.Task, 0)

			if err := tx.
				Where("status != ?", domain.TaskCompleted).
				Find(&tasksDB).Error; err != nil {

				return fmt.Errorf("error get not completed tasks ids: %w", err)
			}

			for _, task := range tasksDB {
				newWorker := workersDB[rand.Intn(len(workersDB))]
				task.UserID = newWorker.ID

				if err := tx.Model(&domain.Task{}).
					Where("id = ? AND status != ?", task.ID, domain.TaskCompleted).
					Updates(&task).Error; err != nil {

					return fmt.Errorf("error assign random worker user to task: %w", err)
				}

				updatedTasks[&task] = &newWorker
			}

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("error random assign users to tasks: %w", err)
	}

	return updatedTasks, err
}

/*
--------------------------- Users ---------------------------
*/

func (i *Instance) CreateOrUpdateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	result := i.database.
		WithContext(
			ctx,
		).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "id",
					},
				},
				DoUpdates: clause.AssignmentColumns(
					[]string{
						"first_name",
						"last_name",
						"username",
						"email",
						"role",
						"updated_at",
					},
				),
			},
		).
		Create(
			user,
		)

	if result.Error != nil {
		logrus.Errorf("r.db.WithContext(ctx).Clauses(...).Create(%v): %v", user, result.Error)

		return uuid.Nil, fmt.Errorf("r.db.WithContext(ctx).Clauses(...).Create(%v): %w", user, result.Error)
	}

	log.Debugf("success create or update user %v in database!", user)

	return user.ID, nil
}

func (i *Instance) GetUsersByRole(ctx context.Context, role domain.UserRoleType) ([]*domain.User, error) {
	var users []*domain.User

	if err := i.database.WithContext(ctx).
		Where("role = ?", role).Find(&users).Error; err != nil {
		log.Errorf(
			"i.database.WithContext(ctx).Where(role = %s).Find(&users): %s",
			role,
			err.Error(),
		)

		return nil, fmt.Errorf("error get users by role %q: %w", role, err)
	}

	log.Debugf("success get users by role %q from database!", role)

	return users, nil
}
