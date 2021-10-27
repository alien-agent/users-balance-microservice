package transaction

import (
	"context"
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
	GetForUser(ctx context.Context, req GetHistoryRequest) ([]entity.Transaction, error)
	CreateOneway(ctx context.Context, req CreateOnewayTransactionRequest) (Transaction, error)
	CreateTwoway(ctx context.Context, req CreateTwowayTransactionRequest) (Transaction, error)
}

// Transaction represents the data about an album.
type Transaction struct {
	entity.Transaction
}

// CreateTwowayTransactionRequest represents a balance update request.
type CreateTwowayTransactionRequest struct {
	SenderId    string `json:"sender_id"` // SenderId is parsed into string instead of UUID to first validate it
	RecipientId string `json:"recipient_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// Validate validates the CreateTwowayTransactionRequest fields.
func (r CreateTwowayTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SenderId, validation.Required, is.UUID),
		validation.Field(&r.RecipientId, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required, validation.Min(0).Exclusive()),
		validation.Field(&r.Description, validation.Length(0, 100)),
	)
}

// CreateOnewayTransactionRequest represents a request to createTwoway a transaction with one participant.
type CreateOnewayTransactionRequest struct {
	OwnerId     string `json:"owner_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description,omitempty"`
}

func (r CreateOnewayTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required),
		validation.Field(&r.Description, validation.Length(0, 100)),
	)
}

type GetHistoryRequest struct {
	OwnerId        string `json:"owner_id"`
	Offset         int       `json:"offset,omitempty"`
	Limit          int       `json:"limit,omitempty"`
	OrderBy        string    `json:"order_by,omitempty"`
	OrderDirection string    `json:"order_direction,omitempty"`
}

func (r GetHistoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.Limit, validation.Min(1)),
		validation.Field(&r.OrderBy, validation.In("transaction_date", "amount")),
		validation.Field(&r.OrderDirection, validation.In("ASC", "DESC")),
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

	return s.depositRepo.Update(ctx, dep)
}

// CreateOneway creates a transaction with only one participant.
// It may be used to reflect a deposit top-up or withdrawal.
func (s service) CreateOneway(ctx context.Context, req CreateOnewayTransactionRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	ownerUUID := uuid.MustParse(req.OwnerId) // req.OwnerId is indeed a valid UUID (because of req.Validate())
	if err := s.modifyBalance(ctx, ownerUUID, req.Amount); err != nil {
		return Transaction{}, err
	}

	transaction := entity.Transaction{
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}
	if req.Amount < 0 {
		transaction.SenderId = ownerUUID
		transaction.Amount = -req.Amount
	} else {
		transaction.RecipientId = ownerUUID
		transaction.Amount = req.Amount
	}

	err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// CreateTwoway creates a transaction with two participants: sender and recipient.
// It may be used to reflect a money transfer from one user to another.
func (s service) CreateTwoway(ctx context.Context, req CreateTwowayTransactionRequest) (Transaction, error) {
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

	transaction := entity.Transaction{
		Id:              0, // will be auto-incremented
		SenderId:        senderUUID,
		RecipientId:     recipientUUID,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: time.Now().UTC(),
	}
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

func (s service) GetForUser(ctx context.Context, req GetHistoryRequest) ([]entity.Transaction, error) {
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
