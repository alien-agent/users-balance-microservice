package transaction

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/dbcontext"
	"users-balance-microservice/pkg/log"
)

// Repository encapsulates the logic to access transactions from the database.
type Repository interface {
	// Create saves a new Transaction in the storage.
	// Transaction tx is assigned an id from database in case of successful transaction.
	Create(ctx context.Context, tx *entity.Transaction) error
	Count(ctx context.Context) (int64, error)
	// GetForUser returns a list of all transactions related to given userId.
	GetForUser(ctx context.Context, ownerId uuid.UUID, orderBy, orderDirection string, offset, limit int) ([]entity.Transaction, error)
}

// repository persists Transaction in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new Transaction repository.
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Create saves a new Transaction record in the database.
// Transaction is assigned an auto-incremented id from database.
func (r repository) Create(ctx context.Context, tx *entity.Transaction) error {
	return r.db.With(ctx).Model(tx).Insert()
}

func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.With(ctx).Select("COUNT(*)").From("transaction").Row(&count)
	return count, err
}

// GetForUser returns all transactions from and to the user with given id.
func (r repository) GetForUser(ctx context.Context, ownerId uuid.UUID, orderBy, orderDirection string, offset, limit int) ([]entity.Transaction, error) {
	var result []entity.Transaction
	query := r.db.With(ctx).Select().
		Where(dbx.Or(dbx.HashExp{"sender_id": ownerId}, dbx.HashExp{"recipient_id": ownerId})).
		Offset(int64(offset)).
		Limit(int64(limit))

	if orderBy != "" {
		if orderDirection == "" {
			query.OrderBy(orderBy)
		} else {
			query.OrderBy(orderBy + " " + orderDirection)
		}
	}

	err := query.All(&result)
	return result, err
}
