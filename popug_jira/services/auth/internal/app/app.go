package app

import (
	"context"
	"fmt"
	"net"
	"time"

	authSrv "github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/auth"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/interceptors"
	usersSrv "github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/users"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/app/users/events"
	"github.com/lyckety/async_arch/popug_jira/services/auth/internal/db/domain"
	jwt "github.com/lyckety/async_arch/popug_jira/services/auth/internal/jwt"
	pbV1Auth "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/api/grpc/auth/v1"
	pbV1Users "github.com/lyckety/async_arch/popug_jira/services/auth/pkg/api/grpc/users/v1"
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
	cudEventer := events.New(cfg.Brokers, cfg.MBTopic, cfg.MBPartition)

	jwtManager := jwt.New(cfg.JWTTokenExpired)

	usersService := usersSrv.New(db, cudEventer)
	authService := authSrv.New(db, jwtManager)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.JWTUnary(jwtManager)),
	)

	pbV1Users.RegisterUsersServiceServer(server, usersService)
	pbV1Auth.RegisterAuthServiceServer(server, authService)

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
