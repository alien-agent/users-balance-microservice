package deposit

import (
	"context"
	"time"

	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/internal/transaction"
	"users-balance-microservice/pkg/log"
)

// Service encapsulates usecase logic for deposits.
type Service interface {
	Get(ctx context.Context, req GetBalanceRequest) (Deposit, error)
	Update(ctx context.Context, req UpdateBalanceRequest) (Transaction, error)
	Transfer(ctx context.Context, req TransferRequest) (Transaction, error)
	GetHistory(ctx context.Context, req GetHistoryRequest) ([]entity.Transaction, error)
	Count(ctx context.Context) (int, error)
}

// Deposit represents the data about a deposit.
type Deposit struct {
	entity.Deposit
}

// Transaction represents the data about a transaction.
type Transaction struct {
	entity.Transaction
}

type service struct {
	depositRepo     Repository
	transactionRepo transaction.Repository
	logger          log.Logger
}

// NewService creates a new album service.
func NewService(depositRepo Repository, transactionRepo transaction.Repository, logger log.Logger) Service {
	return service{depositRepo, transactionRepo, logger}
}

func (s service) modifyBalance(ctx context.Context, ownerId uuid.UUID, amount int64) error {
	dep, err := s.depositRepo.Get(ctx, ownerId)
	if err != nil {
		return err
	}

	dep.Balance += amount
	if dep.Balance < 0 {
		return errors.BadRequest("Insufficient funds to perform operation.")
	}

	return s.depositRepo.Update(ctx, dep)
}

// Get returns the Deposit whose owner whose OwnerId is equal to GetBalanceRequest.OwnerId.
func (s service) Get(ctx context.Context, req GetBalanceRequest) (Deposit, error) {
	if err := req.Validate(); err != nil {
		return Deposit{}, err
	}
	deposit, err := s.depositRepo.Get(ctx, uuid.MustParse(req.OwnerId))
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{deposit}, nil
}

// Create creates a new Deposit.
func (s service) Create(ctx context.Context, ownerId uuid.UUID) (Deposit, error) {
	newDeposit := entity.Deposit{
		OwnerId: ownerId,
		Balance: 0,
	}
	err := s.depositRepo.Create(ctx, newDeposit)
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{newDeposit}, nil
}

// Update changes the balance of Deposit according to UpdateBalanceRequest.
// This method also creates a Transaction record in the database.
// It returns the Transaction which reflects the corresponding balance change in case of success.
func (s service) Update(ctx context.Context, req UpdateBalanceRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	ownerUUID := uuid.MustParse(req.OwnerId) // req.OwnerId is indeed a valid UUID (because of req.Validate())
	if err := s.modifyBalance(ctx, ownerUUID, req.Amount); err != nil {
		return Transaction{}, err
	}

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

	// TODO: Wrap Update() in Transactional
	err := s.transactionRepo.Create(ctx, &tx)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{tx}, nil
}

// Transfer sends money from one user to another according to TransferRequest.
// It returns a Transaction which reflects the corresponding money transfer in case of success.
func (s service) Transfer(ctx context.Context, req TransferRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	// req.SenderId and req.RecipientId are indeed valid UUIDs (checked by req.Validate())
	senderUUID, recipientUUID := uuid.MustParse(req.SenderId), uuid.MustParse(req.RecipientId)
	if err := s.modifyBalance(ctx, senderUUID, -req.Amount); err != nil {
		return Transaction{}, err
	}
	if err := s.modifyBalance(ctx, recipientUUID, req.Amount); err != nil {
		return Transaction{}, err
	}

	tx := entity.Transaction{
		Id:              0, // will be auto-incremented
		SenderId:        senderUUID,
		RecipientId:     recipientUUID,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}
	// TODO: Wrap Transfer() in Transactional()
	if err := s.transactionRepo.Create(ctx, &tx); err != nil {
		return Transaction{}, err
	}
	return Transaction{tx}, nil
}

func (s service) GetHistory(ctx context.Context, req GetHistoryRequest) ([]entity.Transaction, error) {
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

	return s.transactionRepo.GetForUser(ctx, ownerUUID, order, req.Offset, req.Limit)
}

// Count returns a number of Deposits in the database.
// Mainly used for testing purposes.
func (s service) Count(ctx context.Context) (int, error) {
	return s.depositRepo.Count(ctx)
}
