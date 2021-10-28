package deposit

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/log"
)

// Service encapsulates usecase logic for deposits.
type Service interface {
	Get(ctx context.Context, req GetBalanceRequest) (Deposit, error)
	Count(ctx context.Context) (int, error)
}

// Deposit represents the data about an album.
type Deposit struct {
	entity.Deposit
}

// GetBalanceRequest represents a request to getBalance balance of specific user.
type GetBalanceRequest struct {
	OwnerId string `json:"owner_id"`
}

// Validate validates the GetBalanceRequest fields.
func (r GetBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the Deposit whose owner has the specified UUID.
func (s service) Get(ctx context.Context, req GetBalanceRequest) (Deposit, error) {
	if err := req.Validate(); err != nil{
		return Deposit{}, err
	}
	deposit, err := s.repo.Get(ctx, uuid.MustParse(req.OwnerId))
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{deposit}, nil
}

// Create creates a new deposit.
func (s service) Create(ctx context.Context, ownerId uuid.UUID) (Deposit, error) {
	newDeposit := entity.Deposit{
		OwnerId: ownerId,
		Balance: 0,
	}
	err := s.repo.Create(ctx, newDeposit)
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{newDeposit}, nil
}

func (s service) Count(ctx context.Context) (int, error){
	return s.repo.Count(ctx)
}
