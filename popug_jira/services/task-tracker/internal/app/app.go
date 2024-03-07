package app

import (
	"context"
	"fmt"
	"net"
	"time"

	tasksEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/events/tasks"
	usersEvents "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/events/users"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/interceptors"
	tasksSrv "github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/app/task-tracker"
	"github.com/lyckety/async_arch/popug_jira/services/task-tracker/internal/db/domain"
	pbV1Tasks "github.com/lyckety/async_arch/popug_jira/services/task-tracker/pkg/grpc/tasktracker/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *Config
	grpcServer *grpc.Server

	usersCUDEventer *usersEvents.UserCUDEventProcessor
}

func New(cfg *Config, db domain.Repository) *App {
	usersCUDEvents := usersEvents.New(
		db,
		cfg.Brokers,
		cfg.TopicCUDUsers,
		cfg.GroupIDCUDUsers,
		cfg.PartitionCUDUsers,
	)

	tasksCUDEvents := tasksEvents.New(
		cfg.Brokers,
		cfg.TopicCUDTasks,
		cfg.PartitionCUDTasks,
	)

	tasksService := tasksSrv.New(db, tasksCUDEvents)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.JWTUnary),
	)

	pbV1Tasks.RegisterTaskTrackerServiceServer(server, tasksService)

	reflection.Register(server)

	return &App{
		cfg:             cfg,
		grpcServer:      server,
		usersCUDEventer: usersCUDEvents,
	}
}

func (a *App) Run(ctx context.Context) error {
	errGr, grCtx := errgroup.WithContext(ctx)

	errGr.Go(
		func() error {
			if err := a.usersCUDEventer.StartProcessMessages(grCtx); err != nil {
				return err
			}

			return nil
		},
	)

	errGr.Go(
		func() error {
			if err := a.startListen(); err != nil {
				return err
			}

			return nil
		},
	)

	errGr.Go(
		func() error {
			<-grCtx.Done()

			if err := a.stopListen(); err != nil {
				return err
			}

			return nil
		},
	)

	if err := errGr.Wait(); err != nil {
		logrus.Error("errGr.Wait(): ", err)
	}

	return nil
}

func (a *App) startListen() error {
	listener, err := net.Listen("tcp", a.cfg.GRPCBindAddress)
	if err != nil {
		return fmt.Errorf(
			"failed create grpc listener: %q: %w",
			a.cfg.GRPCBindAddress,
			err,
		)
	}

	logrus.Info("gRPC server is listening on ", a.cfg.GRPCBindAddress)

	if err = a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed start gRPC server: %w", err)
	}

	return nil
}

func (a *App) stopListen() error {
	const gracefulShutdownTimeout = time.Second * 10

	ch := make(chan struct{}, 1)

	go func() {
		a.grpcServer.GracefulStop()

		ch <- struct{}{}
	}()

	select {
	case <-ch:
		logrus.Info("grpc server graceful stop")
	case <-time.After(gracefulShutdownTimeout):
		logrus.Info("grpc server force stop")

		a.grpcServer.Stop()
	}

	return nil
}
