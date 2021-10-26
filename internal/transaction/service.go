package transaction

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"users-balance-microservice/internal/deposit"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/log"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	// Get(ctx context.Context, ownerId string) (Transaction, error)
	GetForUser(ctx context.Context, req GetForUserRequest) ([]entity.Transaction, error)
	Create(ctx context.Context, req CreateTransactionRequest) (Transaction, error)
}

// Transaction represents the data about an album.
type Transaction struct {
	entity.Transaction
}

// CreateTransactionRequest represents a balance update request.
type CreateTransactionRequest struct {
	SenderId    uuid.UUID `json:"sender_id"`
	RecipientId uuid.UUID `json:"recipient_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// Validate validates the CreateTransactionRequest fields.
func (r CreateTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SenderId, is.UUID),
		validation.Field(&r.RecipientId, is.UUID),
		validation.Field(&r.Amount, validation.Required, validation.Min(0).Exclusive()),
		validation.Field(&r.Description, validation.Length(0, 80)),
	)
}

type GetForUserRequest struct {
	OwnerId uuid.UUID `json:"owner_id"`
	Offset  int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	// Column name to order transactions by. If empty string passed, no ordering is performed.
	OrderBy   string `json:"order_by,omitempty"`
	// Ascending indicates the sorting direction: true - ascending, false or no value passed - descending.
	Ascending bool   `json:"ascending"`
}

func (r GetForUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.Limit, validation.Min(1)),
	)
}

type service struct {
	transactionRepo Repository
	depositRepo     deposit.Repository
	logger          log.Logger
}

// NewService creates a new album service.
func NewService(transactionRepo Repository, depositRepo deposit.Repository, logger log.Logger) Service {
	return service{transactionRepo, depositRepo, logger}
}

func (s service) modifyBalance(ctx context.Context, ownerId uuid.UUID, amount int64) error {
	dep, err := s.depositRepo.Get(ctx, ownerId)
	if err != nil {
		return err
	}

	dep.Balance += amount
	if dep.Balance < 0 {
		return fmt.Errorf("not enough funds to perform transaction")
	}

	return s.depositRepo.Update(ctx, dep)
}

// Create creates new Transaction based on CreateTransactionRequest and updates its participants' balance.
func (s service) Create(ctx context.Context, req CreateTransactionRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	if req.SenderId != uuid.Nil {
		if err := s.modifyBalance(ctx, req.SenderId, -req.Amount); err != nil {
			return Transaction{}, err
		}
	}
	if req.RecipientId != uuid.Nil {
		if err := s.modifyBalance(ctx, req.RecipientId, req.Amount); err != nil {
			return Transaction{}, err
		}
	}

	transaction := entity.Transaction{
		Id:              0, // will be auto-incremented
		SenderId:        req.SenderId,
		RecipientId:     req.RecipientId,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

func (s service) GetForUser(ctx context.Context, req GetForUserRequest) ([]entity.Transaction, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	order := ""
	if req.OrderBy != "" {
		order = req.OrderBy
		if req.Ascending {
			order += "ASC"
		} else {
			order += "DESC"
		}
	}

	return s.transactionRepo.GetForUser(ctx, req.OwnerId, order, req.Offset, req.Limit)
}
