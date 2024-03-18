package billing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	prodEvents "github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/events/producer"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BillingManager struct {
	storage        domain.Repository
	producerEvents *prodEvents.TxEventSender

	processedCycles *sync.Map
}

func New(
	strg domain.Repository,
	mbProd *prodEvents.TxEventSender,
) *BillingManager {
	return &BillingManager{
		storage:         strg,
		producerEvents:  mbProd,
		processedCycles: &sync.Map{},
	}
}

func (m *BillingManager) Run(ctx context.Context) {
	logrus.Info("Start billing manager...")

	wg := &sync.WaitGroup{}

	errGr, grCtx := errgroup.WithContext(ctx)

	errGr.Go(
		func() error {
			if err := m.runBillingCyclesCreateLoop(grCtx, wg); err != nil {
				logrus.Errorf("app.runBillingCyclesCreateLoop(grCtx): %s", err.Error())

				return fmt.Errorf("app.runBillingCyclesCreateLoop(grCtx): %w", err)
			}

			return nil

		},
	)

	errGr.Go(
		func() error {
			if err := m.runStartBillingCyclesLoop(grCtx, wg); err != nil {
				logrus.Errorf("app.runStartBillingCyclesLoop(grCtx): %s", err.Error())

				return fmt.Errorf("app.runStartBillingCyclesLoop(grCtx): %w", err)
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
			defer logrus.Infof("Canceled billing manager!")

			logrus.Infof("Billing Manager cancel: wait all canceled billing cycle processes...")

			wg.Wait()

			if err := m.storage.AllStartedBillingCyclesToOpened(context.Background()); err != nil {
				logrus.Errorf("m.storage.AllStartedBillingCyclesToOpened(ctx): %s", err.Error())

				return fmt.Errorf("m.storage.AllStartedBillingCyclesToOpened(ctx): %w", err)
			}

			return nil
		},
	)

	if err := errGr.Wait(); err != nil {
		logrus.Panicf("Biiling Manager: fatal error: %s", err)
	}
}

func (m *BillingManager) runBillingCyclesCreateLoop(ctx context.Context, wg *sync.WaitGroup) error {
	logrus.Info("Start loop for creat opened billing cycles...")

	wg.Add(1)
	defer wg.Done()

	iterationTicker := time.NewTicker(3 * time.Second)
	defer iterationTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("canceled runBillingCyclesCreateLoop(ctx)")

			return nil
		case <-iterationTicker.C:
			iterationTicker.Stop()
			if err := m.storage.CreateOpenedBillingCycleForUsers(context.Background()); err != nil {
				logrus.Errorf("m.storage.runBillingCyclesCreateLoop(...): %s", err.Error())
			}

			iterationTicker.Reset(3 * time.Second)
		}
	}
}

func (m *BillingManager) runStartBillingCyclesLoop(ctx context.Context, wg *sync.WaitGroup) error {
	logrus.Info("Start loop for start opened billing cycles...")

	iterationTicker := time.NewTicker(5 * time.Second)
	defer iterationTicker.Stop()

	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("canceled startBillingCyclesCreateLoop(ctx)")

			return nil
		case <-iterationTicker.C:
			iterationTicker.Stop()

			openedBCs, err := m.storage.GetBillingCyclesByStatus(
				context.Background(),
				domain.BillingCycleOpened,
			)

			startedBCs, err := m.storage.GetBillingCyclesByStatus(
				context.Background(),
				domain.BillingCycleStarted,
			)
			if err != nil {
				logrus.Errorf("m.storage.GetBillingCyclesByStatus(...): %s", err.Error())
			}

			// Мотивация  - если в момент старта сервиса в БД для пользователя есть по какой-то причине started bc,
			// то он тоже будет обработан в процессинге
			allOpenedAndStartedCycles := append(openedBCs, startedBCs...)

			for _, cycle := range allOpenedAndStartedCycles {
				if _, ok := m.processedCycles.Load(cycle.ID); !ok {
					m.processedCycles.Store(cycle.ID, true)

					cycle.Status = domain.BillingCycleStarted

					if err := m.storage.CreateOrUpdateBillingCycleStatus(context.Background(), cycle); err != nil {
						return fmt.Errorf("error open billing cycle %v for users: %w", cycle, err)
					}

					wg.Add(1)

					go func(cycle domain.BillingCycle) {
						defer wg.Done()

						defer m.processedCycles.Delete(cycle.ID)

						if err := m.processBillingCycle(ctx, cycle); err != nil {
							logrus.Errorf("error processing billing cycle %v): %s", cycle, err.Error())

							cycle.Status = domain.BillingCycleOpened

							if err := m.storage.CreateOrUpdateBillingCycleStatus(context.Background(), cycle); err != nil {
								logrus.Errorf(
									"error rollback cycle %q to opened status: %s",
									cycle.ID.String(),
									err.Error(),
								)
							}

							return
						}

						cycle.Status = domain.BillingCycleCompleted

						if err := m.storage.CreateOrUpdateBillingCycleStatus(context.Background(), cycle); err != nil {
							logrus.Errorf(
								"error rollback to opened status cycle %q : %s",
								cycle.ID.String(),
								err.Error(),
							)
						}

						logrus.Infof("succes processing billing cycle %v!", cycle)
					}(cycle)
				}
			}

			iterationTicker.Reset(5 * time.Second)
		}
	}
}

func (m *BillingManager) processBillingCycle(ctx context.Context, bc domain.BillingCycle) error {
	logrus.Infof(
		"Billing cycle process %q has been started for user public id %s...",
		bc.ID.String(),
		bc.UserPublicID.String(),
	)

	bc.Status = domain.BillingCycleStarted

	if err := m.storage.CreateOrUpdateBillingCycleStatus(context.Background(), bc); err != nil {
		return fmt.Errorf("m.storage.CreateOrUpdateBillingCycleStatus(ctx, %v): %w", bc, err)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			// TODO: с временем ошибка. Нужно было бы унифицировать везде использование UTC
			if time.Now().UTC().After(bc.EndDate.UTC()) {
				logrus.Infof("Start payment for billing cycle %v", bc)

				txPayment, err := m.storage.ProcessBillingCycle(context.Background(), bc, m.MockPayoutAndNotifyEmail)
				if err != nil {
					bc.Status = domain.BillingCycleOpened

					if err := m.storage.CreateOrUpdateBillingCycleStatus(ctx, bc); err != nil {
						logrus.Errorf("m.storage.CreateOrUpdateBillingCycleStatus(ctx, %v): %s", bc, err.Error())
					}

					return fmt.Errorf("m.storage.CreateOrUpdateBillingCycleStatus(ctx, %v): %w", bc, err)
				}

				if txPayment == nil {
					logrus.Infof("Processing billing cycle: %v nothing to pay", bc)

					bc.Status = domain.BillingCycleCompleted

					if err := m.storage.CreateOrUpdateBillingCycleStatus(ctx, bc); err != nil {
						return fmt.Errorf("m.storage.CreateOrUpdateBillingCycleStatus(ctx, %v): %w", bc, err)
					}

					return nil
				}

				paymentedEvent := &prodEvents.TxPaymentedV1{
					EventTime:      txPayment.DateTime.Unix(),
					PublicID:       txPayment.PublicID,
					WorkerPublicID: txPayment.UserPublicID,
					Cost:           uint32(txPayment.Credit),
				}

				if err := m.producerEvents.Send(ctx, paymentedEvent); err != nil {
					logrus.Errorf("failed send transaction payment event %v: %s", paymentedEvent, err.Error())
				}

				bc.Status = domain.BillingCycleCompleted

				if err := m.storage.CreateOrUpdateBillingCycleStatus(ctx, bc); err != nil {
					return fmt.Errorf("m.storage.CreateOrUpdateBillingCycleStatus(ctx, %v): %w", bc, err)
				}

				logrus.Infof(
					"Billing cycle process %q has been successfully completed for user public id %s!",
					bc.ID.String(),
					bc.UserPublicID.String(),
				)

				return nil
			}
		}
	}

}

func (m *BillingManager) MockPayoutAndNotifyEmail(ctx context.Context, userPublicID uuid.UUID, paymentAmount int) error {
	userDB, err := m.storage.GetUserByPublicID(ctx, userPublicID)
	if err != nil {
		return fmt.Errorf(
			"get user info by public id %s: %w",
			userPublicID.String(),
			err,
		)
	}

	if userDB.Email == "" {
		return fmt.Errorf(
			"not set email for user %q : send email notify imposible",
			userPublicID.String(),
		)
	}

	logrus.Infof(
		"====== PAYOUT TO USER %q (username %q) PAYMENT %d ===============",
		userDB.PublicID,
		userDB.Username,
		paymentAmount,
	)

	logrus.Infof(
		"********* EMAIL NOTIFY TO USER %q (%q) WITH EMAIL (email %q) ***********",
		userDB.Username,
		userDB.PublicID.String(),
		userDB.Email,
	)

	return nil
}
