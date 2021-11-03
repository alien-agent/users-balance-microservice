package deposit

import (
	"context"

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
	// Update updates the changes to the given Deposit to db.
	Update(ctx context.Context, deposit entity.Deposit) error
	// Count returns the number of Deposit records in the database.
	Count(ctx context.Context) (int64, error)
}

// repository persists Deposit in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new Deposit repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the Deposit with the specified OwnerId from the database.
// If Deposit with specified OwnerId does not exist, it is created with balance=0.
func (r repository) Get(ctx context.Context, ownerId uuid.UUID) (entity.Deposit, error) {
	var deposit entity.Deposit
	err := r.db.With(ctx).Select().Model(ownerId, &deposit)
	return deposit, err
}

// Create saves a new Deposit record in the database.
func (r repository) Create(ctx context.Context, deposit entity.Deposit) error {
	return r.db.With(ctx).Model(&deposit).Insert()
}

// Update saves the changes to the Deposit in the database.
func (r repository) Update(ctx context.Context, deposit entity.Deposit) error {
	return r.db.With(ctx).Model(&deposit).Update()
}

// Count returns the number of Deposit records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.With(ctx).Select("COUNT(*)").From("deposit").Row(&count)
	return count, err
}
