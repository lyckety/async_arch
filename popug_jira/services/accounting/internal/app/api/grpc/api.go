package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/app/interceptors"
	"github.com/lyckety/async_arch/popug_jira/services/accounting/internal/db/domain"
	pbV1Acc "github.com/lyckety/async_arch/popug_jira/services/accounting/pkg/api/grpc/accounting/v1"
	pbV1Tx "github.com/lyckety/async_arch/popug_jira/services/accounting/pkg/api/grpc/transaction/v1"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountingService struct {
	pbV1Acc.UnimplementedAccountingServiceServer

	dbIns domain.Repository
}

func New(
	db domain.Repository,
) *AccountingService {
	return &AccountingService{
		dbIns: db,
	}
}

func (s *AccountingService) ShowCurrentBalance(
	ctx context.Context,
	req *pbV1Acc.ShowCurrentBalanceRequest,
) (*pbV1Acc.ShowCurrentBalanceResponse, error) {
	userID := ctx.Value(interceptors.ContextKeyUserID).(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Errorf("ShowCurrentBalance(...): user id must be uuid type: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"ShowCurrentBalance(...): user id must be uuid: %s", err.Error(),
		)
	}

	balance, err := s.dbIns.GetUserBalance(ctx, userUUID)
	if err != nil {
		log.Errorf("s.dbIns.GetUserBalance(ctx, %s): %s", userUUID.String(), err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.GetUserBalance(ctx, %s): %s", userUUID.String(), err.Error(),
		)
	}

	return &pbV1Acc.ShowCurrentBalanceResponse{
		Balance: int64(balance),
	}, nil
}

func (s *AccountingService) ShowTransactionLog(
	ctx context.Context,
	req *pbV1Acc.ShowTransactionLogRequest,
) (*pbV1Acc.ShowTransactionLogResponse, error) {
	userID := ctx.Value(interceptors.ContextKeyUserID).(string)

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Errorf("ShowTransactionLog(...): user id must be uuid type: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"ShowTransactionLog(...): user id must be uuid: %s", err.Error(),
		)
	}

	startDate := time.Unix(req.GetFilter().GetByDateTime().GetStartFrom(), 0)
	var endDate time.Time

	if req.GetFilter().GetByDateTime().GetEndTo() == 0 {
		endDate = time.Now()
	}

	dbTxType := convertRPCTxTypeToDB(req.GetFilter().GetByType())

	txsDB, err := s.dbIns.GetTransactions(
		ctx,
		userUUID,
		startDate, endDate,
		dbTxType,
	)
	if err != nil {
		log.Errorf("s.dbIns.GetTransactions(...): user id must be uuid type: %s", err.Error())

		return nil, status.Errorf(
			codes.Internal,
			"s.dbIns.GetTransactions(...): user id must be uuid type: %s", err.Error(),
		)
	}

	pbTxs := make([]*pbV1Tx.Transaction, len(txsDB))

	for i, tx := range txsDB {
		newPBTx := &pbV1Tx.Transaction{
			PublicId:        tx.PublicID.String(),
			UserId:          tx.UserPublicID.String(),
			Type:            convertDBTxTypeToRPC(tx.Type),
			Debit:           uint32(tx.Debit),
			Credit:          uint32(tx.Credit),
			TaskDescription: tx.TaskDescription,
			Time:            tx.DateTime.Unix(),
		}

		pbTxs[i] = newPBTx
	}

	return &pbV1Acc.ShowTransactionLogResponse{
		Transactions: pbTxs,
	}, nil
}

func convertDBTxTypeToRPC(dbType domain.TransactionType) pbV1Tx.TransactionType {
	switch dbType {
	case domain.TxTypeAssigned:
		return pbV1Tx.TransactionType_TRANSACTION_TYPE_ASSIGNED
	case domain.TxTypeCompleted:
		return pbV1Tx.TransactionType_TRANSACTION_TYPE_COMPLETED
	case domain.TxTypePaymented:
		return pbV1Tx.TransactionType_TRANSACTION_TYPE_PAYMENT
	case domain.TxTypeUnspecified:
		fallthrough
	default:
		return pbV1Tx.TransactionType_TRANSACTION_TYPE_UNKNOWN
	}
}

func convertRPCTxTypeToDB(rpcType pbV1Tx.TransactionType) domain.TransactionType {
	switch rpcType {
	case pbV1Tx.TransactionType_TRANSACTION_TYPE_ASSIGNED:
		return domain.TxTypeAssigned
	case pbV1Tx.TransactionType_TRANSACTION_TYPE_COMPLETED:
		return domain.TxTypeCompleted
	case pbV1Tx.TransactionType_TRANSACTION_TYPE_PAYMENT:
		return domain.TxTypePaymented
	case pbV1Tx.TransactionType_TRANSACTION_TYPE_UNSPECIFIED:
		fallthrough
	default:
		return domain.TxTypeUnspecified
	}
}
