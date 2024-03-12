package app

import (
	"context"
	"fmt"
	"net"
	"time"

	accAPI "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/api/grpc"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/billing"
	consEvents "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/events/consumer"
	prodEvents "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/events/producer"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/interceptors"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	pbV1Acc "github.com/lyckety/async_arch/popug_jira/services/accounting/pkg/api/grpc/accounting/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *Config
	grpcServer *grpc.Server

	billingManager *billing.BillingManager

	prodTxEvents *prodEvents.TxEventSender

	usersCUDEventer *consEvents.UsersCUDProcessor
	tasksCUDEventer *consEvents.TasksStreamProcessor
	tasksBEventer   *consEvents.TasksBEProcessor
}

func New(cfg *Config, db domain.Repository) *App {
	txEventsProducer := prodEvents.New(cfg.Brokers, cfg.TopicCUDTxs, cfg.PartitionCUDTasks)

	billing := billing.New(db, txEventsProducer)

	usersCUDEvents := consEvents.NewUsersCUDProcessor(
		db,
		cfg.Brokers,
		cfg.TopicCUDUsers,
		cfg.EventsGroupID,
		cfg.PartitionCUDUsers,
	)

	tasksCUDEvents := consEvents.NewTasksStreamProcessor(
		db,
		cfg.Brokers,
		cfg.TopicCUDTasks,
		cfg.EventsGroupID,
		cfg.PartitionCUDTasks,
	)

	tasksBE := consEvents.NewTasksBEProcessor(
		db,
		cfg.Brokers,
		cfg.TopicBETasks,
		cfg.EventsGroupID,
		cfg.PartitionBETasks,
		txEventsProducer,
	)

	accService := accAPI.New(db)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.JWTUnary),
	)

	pbV1Acc.RegisterAccountingServiceServer(server, accService)

	reflection.Register(server)

	return &App{
		cfg:             cfg,
		grpcServer:      server,
		billingManager:  billing,
		usersCUDEventer: usersCUDEvents,
		tasksCUDEventer: tasksCUDEvents,
		tasksBEventer:   tasksBE,
	}
}

func (a *App) Run(ctx context.Context) error {
	errGr, grCtx := errgroup.WithContext(ctx)

	errGr.Go(
		func() error {
			a.billingManager.Run(grCtx)

			return nil
		},
	)

	errGr.Go(
		func() error {
			a.usersCUDEventer.StartProcessMessages(grCtx)

			return nil
		},
	)

	errGr.Go(
		func() error {
			a.tasksCUDEventer.StartProcessMessages(grCtx)

			return nil
		},
	)

	errGr.Go(
		func() error {
			a.tasksBEventer.StartProcessMessages(grCtx)

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
