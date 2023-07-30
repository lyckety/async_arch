package postgresql

import (
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint:nolintlint
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrDirtyDBVersion = errors.New("dirty db version")

type Instance struct {
	config   *Config
	database *gorm.DB

	cacheTableNames []string
}

func New(cfg *Config) *Instance {
	return &Instance{
		config:          cfg,
		cacheTableNames: make([]string, 0),
	}
}

func (i *Instance) Connect() error {
	var err error

	i.database, err = gorm.Open(
		postgres.Open(i.config.GetDSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		logrus.Errorf("Cannot connect to DB. Reason: %s", err.Error())

		return fmt.Errorf("failed connect to database: %w", err)
	}

	logrus.Info("Connected to database")

	return nil
}

func (i *Instance) Disconnect() {
	if sqlDB, _ := i.database.DB(); sqlDB != nil {
		_ = sqlDB.Close()
	}
}

func (i *Instance) IsConnected() bool {
	sqlDB, err := i.database.DB()

	return i.database != nil && sqlDB != nil && err == nil && sqlDB.Ping() == nil
}
