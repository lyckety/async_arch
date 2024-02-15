package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/lyckety/async_arch/golang-service/internal/db/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (i *Instance) CreateOrUpdateUser(ctx context.Context, user *domain.User) error {
	result := i.database.
		WithContext(
			ctx,
		).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "username",
					},
				},
				DoUpdates: clause.AssignmentColumns(
					[]string{
						"email",
						"date_of_birth",
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

		return fmt.Errorf("r.db.WithContext(ctx).Clauses(...).Create(%v): %w", user, result.Error)
	}

	return nil
}

func (i *Instance) DeleteUser(
	ctx context.Context,
	userName string,
) error {
	result := i.database.
		WithContext(
			ctx,
		).
		Where(
			"username = ?", userName).
		Delete(
			&domain.User{})
	if result.Error != nil {
		logrus.Errorf("i.database.WithContext(ctx).Where(name_table = %s)(&domain.User{}): %v", userName, result.Error)

		return fmt.Errorf(
			"i.database.WithContext(ctx).Where(name_table = %s)(&domain.User{}): %w",
			userName,
			result.Error,
		)
	}

	return nil
}

func (i *Instance) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var allUsers []*domain.User

	if err := i.database.
		WithContext(
			ctx,
		).
		Find(
			&allUsers,
		).Error; err != nil {
		logrus.Errorf(
			"i.database.WithContext(ctx).Find(&allUsers): %v",
			err,
		)

		return nil,
			fmt.Errorf(
				"i.database.WithContext(ctx).Find(&allUsers): %w",
				err,
			)
	}

	return allUsers, nil
}

func (i *Instance) GetUserByName(ctx context.Context, userName string) (*domain.User, error) {
	var user domain.User

	if err := i.database.
		WithContext(
			ctx,
		).
		Where(
			"username = ?", userName,
		).
		First(
			&user,
		).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil,
				fmt.Errorf(
					"i.database.WithContext(ctx).Where(%s).First(&user): not found username %s",
					userName,
					userName,
				)
		} else {
			logrus.Errorf(
				"i.database.WithContext(ctx).Where(%s).First(&user): %v",
				userName,
				err,
			)

			return nil,
				fmt.Errorf(
					"i.database.WithContext(ctx).Where(%s).First(&user): %w",
					userName,
					err,
				)
		}

	}

	return &user, nil
}
