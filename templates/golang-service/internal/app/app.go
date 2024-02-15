package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/lyckety/async_arch/golang-service/internal/app/consumer"
	"github.com/lyckety/async_arch/golang-service/internal/app/producer"
	usersSrv "github.com/lyckety/async_arch/golang-service/internal/app/users"
	"github.com/lyckety/async_arch/golang-service/internal/db/domain"
	pbV1Cons "github.com/lyckety/async_arch/golang-service/pkg/grpc/kafkaconsumer/v1"
	pbV1Prod "github.com/lyckety/async_arch/golang-service/pkg/grpc/kafkaproducer/v1"
	pbV1Users "github.com/lyckety/async_arch/golang-service/pkg/grpc/users/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg        *Config
	grpcServer *grpc.Server
}

func New(cfg *Config, db domain.Repository) *App {
	producerService := producer.New(
		cfg.Brokers,
	)

	consumerService := consumer.New(
		cfg.Brokers,
	)

	usersService := usersSrv.New(db)

	server := grpc.NewServer()

	pbV1Prod.RegisterKafkaProducerServiceServer(server, producerService)
	pbV1Cons.RegisterKafkaConsumerServiceServer(server, consumerService)

	pbV1Users.RegisterUsersServiceServer(server, usersService)

	reflection.Register(server)

	return &App{
		cfg:        cfg,
		grpcServer: server,
	}
}

func (a *App) Run(ctx context.Context) error {
	errGr, grCtx := errgroup.WithContext(ctx)

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
			select {
			case <-ctx.Done():
			case <-grCtx.Done():
			}

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
