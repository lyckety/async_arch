package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/config"
	db "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/postgresql"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logger.Panic("failed init config: %w", err)
	}

	initLogger(cfg.GetLogLevel())

	database := db.New(
		&db.Config{
			DBHost:         cfg.GetDBHost(),
			DBPort:         cfg.GetDBPort(),
			DBName:         cfg.GetDBName(),
			DBUsername:     cfg.GetDBUserName(),
			DBPassword:     cfg.GetDBPassword(),
			MigrationsPath: cfg.GetMigrationsPath(),
		},
	)

	if err := database.Connect(); err != nil {
		logrus.Panic(err)
	}
	defer database.Disconnect()

	if err := database.Migrate(); err != nil {
		logrus.Panicf("failed migrate database: %s", err)
	}

	app := app.New(
		&app.Config{
			GRPCBindAddress:   cfg.GetGRPCBindAddress(),
			Brokers:           cfg.GetBrokers(),
			TopicCUDUsers:     cfg.GetUsersCUDTopicName(),
			GroupIDCUDUsers:   cfg.GetUsersCUDGroupID(),
			PartitionCUDUsers: cfg.GetUsersCUDPartition(),
			TopicCUDTasks:     cfg.GetTasksCUDTopicName(),
			PartitionCUDTasks: cfg.GetTasksCUDPartition(),
		},
		database,
	)

	ctx, cancel := context.WithCancel(context.Background())
	errGr, grCtx := errgroup.WithContext(ctx)

	errGr.Go(
		func() error {
			if err := app.Run(grCtx); err != nil {
				return fmt.Errorf("app.Run(): %w", err)
			}

			return nil

		},
	)

	errGr.Go(
		func() error {
			waitForSignal(grCtx, cancel)

			return nil
		},
	)

	if err := errGr.Wait(); err != nil {
		logrus.Panicf("fatal error: %s", err)
	}

}

func waitForSignal(ctx context.Context, cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigs:
		logrus.Infof("catched %v\n", sig)

		cancel()
	case <-ctx.Done():
		logrus.Trace("context Done")
	}
}

func initLogger(level logger.Level) {
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)

	logger.SetLevel(level)

	logger.SetFormatter(
		&logger.TextFormatter{ //nolint:exhaustruct
			FullTimestamp: true,
		},
	)
}

func waitForStopSignal(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigs:
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				return
			}
		}
	}
}
