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
	Get(ctx context.Context, ownerId uuid.UUID) (Deposit, error)
	Update(ctx context.Context, request UpdateBalanceRequest) (Deposit, error)
}

// Deposit represents the data about an album.
type Deposit struct {
	entity.Deposit
}

// GetBalanceRequest represents a request to get balance of specific user.
type GetBalanceRequest struct {
	OwnerId uuid.UUID `json:"owner_id"`
}

// Validate validates the GetBalanceRequest fields.
func (r GetBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
	)
}

// UpdateBalanceRequest represents a balance update request.
type UpdateBalanceRequest struct {
	OwnerId     uuid.UUID `json:"owner_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// Validate validates the GetBalanceRequest fields.
func (r UpdateBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required),
		validation.Field(&r.Description, validation.Length(0, 80)),
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
func (s service) Get(ctx context.Context, ownerId uuid.UUID) (Deposit, error) {
	deposit, err := s.repo.Get(ctx, ownerId)
	if err != nil {
		return Deposit{}, err
	}
	return Deposit{deposit}, nil
}

// Create creates a new album.
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

// Update updates the album with the specified ID.
func (s service) Update(ctx context.Context, req UpdateBalanceRequest) (Deposit, error) {
	if err := req.Validate(); err != nil {
		return Deposit{}, err
	}

	/*
		album, err := s.Get(ctx, id)
		if err != nil {
			return album, err
		}
		album.Name = req.Name
		album.UpdatedAt = time.Now()

		if err := s.repo.Create(ctx, album.Deposit); err != nil {
			return album, err
		}
		return album, nil*/
	return Deposit{}, nil
}
