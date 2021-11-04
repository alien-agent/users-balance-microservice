package transaction

import (
	"context"
	"time"

	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/requests"
	"users-balance-microservice/pkg/log"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	// CreateUpdateTransaction creates a Transaction based on UpdateBalanceRequest.
	CreateUpdateTransaction(ctx context.Context, req requests.UpdateBalanceRequest) (Transaction, error)
	// CreateTransferTransaction creates a Transaction based on TransferRequest.
	CreateTransferTransaction(ctx context.Context, req requests.TransferRequest) (Transaction, error)
	// GetHistory returns a list of all transactions related to the user with the given ID.
	GetHistory(ctx context.Context, req requests.GetHistoryRequest) ([]entity.Transaction, error)
	// Count returns a number of all Transactions in the database. Mainly used for testing purposes.
	Count(ctx context.Context) (int64, error)
}

// Transaction represents the data about a transaction.
type Transaction struct {
	entity.Transaction
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new Transaction service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) CreateUpdateTransaction(ctx context.Context, req requests.UpdateBalanceRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	ownerUUID := uuid.MustParse(req.OwnerId)
	tx := entity.Transaction{
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}
	if req.Amount < 0 {
		tx.SenderId = ownerUUID
		tx.Amount = -req.Amount
	} else {
		tx.RecipientId = ownerUUID
		tx.Amount = req.Amount
	}

	err := s.repo.Create(ctx, &tx)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{tx}, err
}

func (s service) CreateTransferTransaction(ctx context.Context, req requests.TransferRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	senderUUID, recipientUUID := uuid.MustParse(req.SenderId), uuid.MustParse(req.RecipientId)
	tx := entity.Transaction{
		Id:              0, // will be auto-incremented
		SenderId:        senderUUID,
		RecipientId:     recipientUUID,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}

	err := s.repo.Create(ctx, &tx)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{tx}, err
}

func (s service) GetHistory(ctx context.Context, req requests.GetHistoryRequest) ([]entity.Transaction, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ownerUUID := uuid.MustParse(req.OwnerId)
	order := ""
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.OrderDirection != "" {
			order = order + " " + req.OrderDirection
		}
	}

	return s.repo.GetForUser(ctx, ownerUUID, order, req.Offset, req.Limit)
}

func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
