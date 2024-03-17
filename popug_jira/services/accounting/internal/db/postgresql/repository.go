package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
--------------------------- Tasks ---------------------------
*/

func (i *Instance) CreateOrUpdateTask(ctx context.Context, task *domain.Task) (uuid.UUID, error) {
	updatedColumns := []string{
		"updated_at",
	}

	if strings.ReplaceAll(task.Description, " ", "") != "" {
		updatedColumns = append(updatedColumns, "description")
	}

	if strings.ReplaceAll(task.JiraId, " ", "") != "" {
		updatedColumns = append(updatedColumns, "jira_id")
	}

	err := i.database.
		WithContext(
			ctx,
		).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "public_id",
					},
				},
				DoUpdates: clause.AssignmentColumns(
					updatedColumns,
				),
			},
		).
		Create(
			task,
		).Error

	if err != nil {
		return uuid.Nil, fmt.Errorf("error create or update task %v: %w", task, err)
	}

	log.Debugf("success create or update task %v in database!", task)

	return task.PublicID, nil
}

func (i *Instance) GetTaskByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Task, error) {
	var task *domain.Task

	if err := i.database.WithContext(ctx).
		Where("public_id = ?", publicID.String()).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warningf("not found task by public id %q", publicID.String())

			return nil, nil
		}

		return nil, fmt.Errorf("error get task by public id %q: %w", publicID.String(), err)
	}

	log.Debugf("success get task by public id %q from database!", publicID.String())

	return task, nil
}

/*
--------------------------- Users ---------------------------
*/

func (i *Instance) CreateOrUpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := i.database.
		WithContext(
			ctx,
		).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "public_id",
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
		return nil, fmt.Errorf("r.db.WithContext(ctx).Clauses(...).Create(%v): %w", user, result.Error)
	}

	log.Debugf("success create or update user %v in database!", user)

	return user, nil
}

func (i *Instance) GetUserByPublicID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User

	if err := i.database.WithContext(ctx).
		Where("public_id = ?", id.String()).
		Find(&user).Error; err != nil {

		return nil, fmt.Errorf("error get user by id %q: %w", id.String(), err)
	}

	log.Debugf("success get user by id %q from database!", id.String())

	return user, nil
}

/*
--------------------------- Accounts ---------------------------
*/

func (i *Instance) GetUserBalance(ctx context.Context, publicID uuid.UUID) (int, error) {
	var user domain.User

	if err := i.database.WithContext(ctx).
		Select("balance").
		First(&user, "public_id = ?", publicID).Error; err != nil {

		return 0, fmt.Errorf(
			"error get balance by user public_id %s: %w",
			publicID.String(),
			err,
		)
	}

	return user.Balance, nil
}

func (i *Instance) GetTransactions(
	ctx context.Context,
	userPublicID uuid.UUID,
	startDate time.Time, endDate time.Time,
	transactionType domain.TransactionType,
) ([]domain.TransactionWithDescriptionTask, error) {
	var tx domain.Transaction
	var taskDB domain.Task

	var transactions []domain.TransactionWithDescriptionTask

	query := i.database.WithContext(ctx).
		Table(tx.TableName()).
		Select(
			fmt.Sprintf("%s.*, %s.description as task_description, %s.jira_id as jira_id",
				tx.TableName(),
				taskDB.TableName(),
				taskDB.TableName(),
			),
		).
		Joins(
			fmt.Sprintf("left join tasks on %s.task_public_id = %s.public_id", tx.TableName(), taskDB.TableName()),
		).
		Where(
			fmt.Sprintf(
				"%s.user_public_id = ? AND %s.created_at BETWEEN ? AND ?",
				tx.TableName(),
				tx.TableName(),
			),
			userPublicID,
			startDate,
			endDate,
		)

	if userPublicID != uuid.Nil {
		query = query.Where(fmt.Sprintf("%s.user_public_id = ?", tx.TableName()), userPublicID)
	}

	if transactionType != domain.TxTypeUnspecified {
		query = query.Where(fmt.Sprintf("%s.type = ?", tx.TableName()), transactionType)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

/*
--------------------------- Transactions ---------------------------
*/

func (i *Instance) CreateDebitTransaction(
	ctx context.Context,
	taskPublicID uuid.UUID,
	userPublicID uuid.UUID,
	timeStamp time.Time,
) (*domain.Transaction, error) {
	userDB := &domain.User{}
	taskDB := &domain.Task{}

	tx := i.database.Begin()

	if err := retryFunc(
		5,
		time.Second*5,
		func() error {
			if err := i.database.WithContext(ctx).Where("public_id = ?", taskPublicID).First(taskDB).Error; err != nil {
				return fmt.Errorf(
					"tx.Where('public_id = ?', %s).First(%v): %s",
					taskPublicID.String(),
					taskDB,
					err.Error(),
				)
			}

			return nil
		},
	); err != nil {
		return nil, fmt.Errorf(
			"error find for transaction related task with public id %q: %s",
			taskPublicID.String(),
			err.Error(),
		)
	}

	txDB := &domain.Transaction{
		UserPublicID: userPublicID,
		TaskPublicID: taskPublicID,
		Type:         domain.TxTypeCompleted,
		Debit:        taskDB.CostComplete,
		DateTime:     timeStamp,
	}
	if err := tx.Create(txDB).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf(
			"create debit transaction(ctx, %s, %s): %s",
			taskPublicID.String(),
			userPublicID.String(),
			err.Error(),
		)
	}
	userDB.PublicID = userPublicID

	if err := tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "public_id",
				},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"balance":    gorm.Expr("users.balance + ?", txDB.Debit),
				"updated_at": time.Now(),
			}),
		},
	).
		Create(
			userDB,
		).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf(
			"error create debit transaction: save db transaction (%s, %s): %w",
			taskPublicID.String(),
			userPublicID.String(),
			err,
		)
	}

	tx.Commit()

	return txDB, nil
}

func (i *Instance) CreateCreditTransaction(
	ctx context.Context,
	taskPublicID uuid.UUID,
	userPublicID uuid.UUID,
	timeStamp time.Time,
) (*domain.Transaction, error) {
	userDB := &domain.User{}
	taskDB := &domain.Task{}
	tx := i.database.WithContext(ctx).Begin()

	if err := retryFunc(
		5,
		time.Second*5,
		func() error {
			if err := i.database.WithContext(ctx).
				Where("public_id = ?", taskPublicID).
				First(taskDB).Error; err != nil {
				return fmt.Errorf("error create credit transaction: find task related task: %w", err)
			}

			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("error create credit transaction: find task related task: %w", err)
	}

	txDB := &domain.Transaction{
		UserPublicID: userPublicID,
		TaskPublicID: taskPublicID,
		Type:         domain.TxTypeAssigned,
		Credit:       taskDB.CostAssign,
		DateTime:     timeStamp,
	}
	if err := tx.Create(txDB).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("error create credit transaction: %w", err)
	}

	userDB.PublicID = userPublicID
	userDB.Balance -= txDB.Credit

	if err := tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "public_id",
				},
			},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"balance":    gorm.Expr("users.balance - ?", txDB.Credit),
					"updated_at": time.Now(),
				},
			),
		},
	).
		Create(
			userDB,
		).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf(
			"error create credit transaction (task: %q, user: %q): save db transaction:: %w",
			taskPublicID.String(),
			userPublicID.String(),
			err,
		)
	}

	tx.Commit()

	return txDB, nil
}

func (i *Instance) CreateOrUpdateBillingCycleStatus(ctx context.Context, bc domain.BillingCycle) error {
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
						"status",
						"updated_at",
					},
				),
			},
		).
		Create(
			&bc,
		)

	if result.Error != nil {
		return fmt.Errorf("error update billing cycle  %v status: %w", bc, result.Error)
	}

	log.Debugf("success create or update bc %v in database!", bc)

	return nil
}

func (i *Instance) CreateOpenedBillingCycleForUsers(ctx context.Context) error {
	tx := i.database.WithContext(ctx).Begin()

	var users []domain.User
	if err := tx.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		var billingCycle domain.BillingCycle

		// // Включить для выполнения условия задачи -  в конце дня выполняется выплата положительного баланса
		// now := time.Now()
		// startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, now.Minute()+5, 0, 0, now.Location())

		startDate := time.Now().UTC()
		endDate := startDate.Add(3 * time.Minute).UTC()

		if err := tx.Where("user_public_id = ? AND status IN (?, ?)", user.PublicID, domain.BillingCycleOpened, domain.BillingCycleStarted).
			First(&billingCycle).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				billingCycle = domain.BillingCycle{
					UserPublicID: user.PublicID,
					StartDate:    startDate,
					EndDate:      endDate,
					Status:       domain.BillingCycleOpened,
				}
				if err := tx.Create(&billingCycle).Error; err != nil {
					logrus.Errorf(
						"error create billing cycle for user_public_id %s: %s", user.PublicID.String(), err.Error(),
					)

					continue
				}

				log.Debugf("success create opened bc %v in database!", billingCycle)
			}
		}
	}

	tx.Commit()

	return nil
}

func (i *Instance) GetBillingCyclesByStatus(
	ctx context.Context,
	status domain.BillingCycleStatus,
) ([]domain.BillingCycle, error) {
	var billingCycles []domain.BillingCycle

	if err := i.database.WithContext(ctx).
		Where("status = ?", status).Find(&billingCycles).Error; err != nil {

		return billingCycles,
			fmt.Errorf(
				"error get billing cycles by status %s: %w", status, err,
			)
	}

	return billingCycles, nil
}

func (i *Instance) ProcessBillingCycle(
	ctx context.Context,
	bc domain.BillingCycle,
	payoutFunc domain.HandleFuncPayoutToUser,
) (*domain.Transaction, error) {
	tx := i.database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &domain.User{
		PublicID: bc.UserPublicID,
	}

	if err := retryFunc(
		5,
		time.Second*5,
		func() error {
			if err := i.database.WithContext(ctx).
				Select("balance").
				First(&user, "public_id = ?", bc.UserPublicID).
				Error; err != nil {
				return fmt.Errorf("error get balance by user public_id %s: %w", bc.UserPublicID.String(), err)
			}

			return nil
		},
	); err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("error get user balance for calculate payment: %w", err)
	}

	payment := user.Balance

	if payment <= 0 {
		tx.Commit()

		logrus.Infof(
			"process billing cycle %q for user %q in database is finished (current payment amount /<= 0)",
			bc.ID.String(),
			bc.UserPublicID.String(),
		)

		return nil, nil
	}

	resetedBalance := map[string]interface{}{
		"balance":    0,
		"updated_at": time.Now(),
	}

	if err := tx.Model(&user).
		Where("public_id = ?", bc.UserPublicID).
		Updates(resetedBalance).Error; err != nil {
		tx.Rollback()

		return nil,
			fmt.Errorf(
				"error update balance in billing cycle process (payment == %d): %s",
				payment,
				err.Error(),
			)
	}

	if err := payoutFunc(ctx, user.PublicID, payment); err != nil {
		tx.Rollback()

		return nil, err
	}

	txDB := &domain.Transaction{
		UserPublicID: user.PublicID,
		Type:         domain.TxTypePaymented,
		Credit:       payment,
		DateTime:     time.Now(),
	}
	if err := tx.Create(txDB).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("error create payment transaction: %w", err)
	}

	tx.Commit()

	logrus.Infof(
		"process billing cycle in database is finished success (sum %d => user %q)",
		payment,
		bc.UserPublicID,
	)

	return txDB, nil
}

func (i *Instance) AllStartedBillingCyclesToOpened(
	ctx context.Context,
) error {
	if err := i.database.WithContext(ctx).Model(&domain.BillingCycle{}).
		Where("status = ?", domain.BillingCycleStarted).
		Update("status", domain.BillingCycleOpened).
		Error; err != nil {

		return fmt.Errorf("error all started billing cycles to opened: %w", err)
	}

	return nil
}
