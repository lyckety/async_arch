package app

import (
	"context"

	consEvents "github.com/lyckety/async_arch/popug_jira/services/analytics/internal/app/events/consumer"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg *Config

	usersCUDEventer *consEvents.UsersCUDProcessor
	tasksCUDEventer *consEvents.TasksStreamProcessor
	txCUDEventer    *consEvents.TxCUDProcessor
}

func New(cfg *Config) *App {
	usersCUDEvents := consEvents.NewUsersCUDProcessor(
		cfg.Brokers,
		cfg.TopicCUDUsers,
		cfg.EventsGroupID,
		cfg.PartitionCUDUsers,
	)

	tasksCUDEvents := consEvents.NewTasksStreamProcessor(
		cfg.Brokers,
		cfg.TopicCUDTasks,
		cfg.EventsGroupID,
		cfg.PartitionCUDTasks,
	)

	txCUDEvents := consEvents.NewTxCUDProcessor(
		cfg.Brokers,
		cfg.TopicCUDTxs,
		cfg.EventsGroupID,
		cfg.PartitionCUDTxs,
	)

	return &App{
		cfg:             cfg,
		usersCUDEventer: usersCUDEvents,
		tasksCUDEventer: tasksCUDEvents,
		txCUDEventer:    txCUDEvents,
	}
}

func (a *App) Run(ctx context.Context) error {
	errGr, grCtx := errgroup.WithContext(ctx)

	errGr.Go(
		func() error {
			logrus.Info("Start users stream processor...")
			a.usersCUDEventer.StartProcessMessages(grCtx)

			logrus.Info("Stopped users stream processor")

			return nil
		},
	)

	errGr.Go(
		func() error {
			logrus.Info("Start tasks stream processor...")

			a.tasksCUDEventer.StartProcessMessages(grCtx)

			logrus.Info("Stopped tasks stream processor")

			return nil
		},
	)

	errGr.Go(
		func() error {
			logrus.Info("Start transactions stream processor...")

			a.txCUDEventer.StartProcessMessages(grCtx)

			logrus.Info("Stopped transactions stream processor")

			return nil
		},
	)

	errGr.Go(
		func() error {
			<-grCtx.Done()

			return nil
		},
	)

	if err := errGr.Wait(); err != nil {
		logrus.Error("errGr.Wait(): ", err)
	}

	return nil
}
