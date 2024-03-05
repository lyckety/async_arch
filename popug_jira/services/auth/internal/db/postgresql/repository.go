package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (i *Instance) GetUserByUsername(ctx context.Context, userName string) (*domain.User, error) {
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
			return nil, nil
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

func (i *Instance) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	result := i.database.WithContext(ctx).Create(&user)

	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("error create user %s: %w", user.Username, result.Error)
	}

	return user.ID, nil
}

func (i *Instance) UpdateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	existedUser := *user

	result := i.database.
		WithContext(ctx).
		Where(
			"username = ?",
			user.Username,
		).
		First(&existedUser)
	if result.Error != nil {
		log.Errorf("error update user %s: %s", user.Username, result.Error)

		return uuid.Nil, fmt.Errorf("error update user %s: %w", user.Username, result.Error)
	}

	if err := i.database.
		Model(user).
		Where(
			"username = ?",
			user.Username,
		).
		Updates(
			user,
		).Error; err != nil {
		log.Errorf("error update user %s: %s", user.Username, result.Error)

		return uuid.Nil, fmt.Errorf("error update user %s: %w", user.Username, result.Error)
	}

	return existedUser.ID, nil
}
