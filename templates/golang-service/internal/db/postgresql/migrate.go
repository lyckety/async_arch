package postgresql

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migpsql "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (i *Instance) Migrate() error {
	gormDB, _ := gorm.Open(postgres.Open(i.config.GetDSN()))
	sqlDB, _ := gormDB.DB()

	driver, err := migpsql.WithInstance(sqlDB, &migpsql.Config{})
	if err != nil {
		return fmt.Errorf("migpsql.WithInstance(...): %w", err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", i.config.getMigrationsPath()), "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance(...): %w", err)
	}
	defer m.Close()

	if err := checkVersion(m); err != nil {
		return fmt.Errorf("checkVersion(...): %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Info("No migrations required")

			return nil
		}

		return fmt.Errorf("m.Up(): %w", err)
	}

	if err := checkVersion(m); err != nil {
		return fmt.Errorf("checkVersion(...): %w", err)
	}

	return nil
}

func checkVersion(m *migrate.Migrate) error {
	version, dirty, err := m.Version()

	logrus.Infof("Current version: %d", version)

	if dirty {
		return fmt.Errorf("%w: %d", ErrDirtyDBVersion, version)
	}

	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("m.Version(): %w", err)
	}

	return nil
}
