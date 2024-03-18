package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lyckety/async_arch/popug_jira/services/analytics/internal/app"
	"github.com/lyckety/async_arch/popug_jira/services/analytics/internal/config"
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

	app := app.New(
		&app.Config{
			Brokers:           cfg.GetBrokers(),
			EventsGroupID:     cfg.GetEventsGroupID(),
			TopicCUDUsers:     cfg.GetUsersCUDTopicName(),
			PartitionCUDUsers: cfg.GetUsersCUDPartition(),
			TopicCUDTasks:     cfg.GetTasksCUDTopicName(),
			PartitionCUDTasks: cfg.GetTasksCUDPartition(),
			TopicCUDTxs:       cfg.GetTxCUDTopicName(),
			PartitionCUDTxs:   cfg.GetTxCUDPartition(),
		},
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
			waitForStopSignal(grCtx, cancel)

			return nil
		},
	)

	if err := errGr.Wait(); err != nil {
		logrus.Panicf("fatal error: %s", err)
	}

}

func waitForStopSignal(ctx context.Context, cancel context.CancelFunc) {
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
		&logger.JSONFormatter{ //nolint:exhaustruct
			PrettyPrint: true,
		},
	)
}
