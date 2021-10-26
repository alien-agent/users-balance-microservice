package deposit

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/dbcontext"
	"users-balance-microservice/pkg/log"
)

// Repository encapsulates the logic to access deposits from the database.
type Repository interface {
	// Get returns the Deposit with the specified owner's UUID.
	Get(ctx context.Context, ownerId uuid.UUID) (entity.Deposit, error)
	// Create saves a new Deposit in the storage.
	Create(ctx context.Context, deposit entity.Deposit) error
	// Update updates the given Deposit to db.
	Update(ctx context.Context, deposit entity.Deposit) error
}

var NegativeBalanceError = errors.New("cannot save deposit with negative balance")

// repository persists Deposit in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new album repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the Deposit with the specified OwnerId from the database.
// If Deposit with specified OwnerId does not exist, it is created with balance=0.
func (r repository) Get(ctx context.Context, ownerId uuid.UUID) (entity.Deposit, error) {
	var deposit entity.Deposit
	err := r.db.With(ctx).Select().Model(ownerId, &deposit)

	// Create if does not exist
	if err == sql.ErrNoRows{
		err = r.Create(ctx, entity.Deposit{
			OwnerId: ownerId,
			Balance: 0,
		})
		if err == nil {
			deposit, err = r.Get(ctx, ownerId)
		}
	}

	return deposit, err
}

// Create saves a new Deposit record in the database.
func (r repository) Create(ctx context.Context, deposit entity.Deposit) error {
	return r.db.With(ctx).Model(&deposit).Insert()
}

// Update saves the changes to the Deposit in the database.
// Deposits with negative balance are rejected.
func (r repository) Update(ctx context.Context, deposit entity.Deposit) error {
	if deposit.Balance < 0{
		return NegativeBalanceError
	}
	return r.db.With(ctx).Model(&deposit).Update()
}