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
	// GetForUser returns a list of all transactions related to given userId.
	GetForUser(ctx context.Context, ownerId uuid.UUID, order string, offset, limit int) ([]entity.Transaction, error)
}

// repository persists Transaction in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new album repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Create saves a new Transaction record in the database.
// Transaction tx is assigned an id from database in case of successful transaction.
func (r repository) Create(ctx context.Context, tx *entity.Transaction) error {
	return r.db.With(ctx).Model(tx).Insert()
}

// GetForUser returns all transactions from and to the user(Transaction) with given id.
func (r repository) GetForUser(ctx context.Context, ownerId uuid.UUID, order string, offset, limit int) ([]entity.Transaction, error) {
	var result []entity.Transaction
	query := r.db.With(ctx).Select().Where(dbx.Or(dbx.HashExp{"sender_id": ownerId}, dbx.HashExp{"recipient_id": ownerId}))

	if offset >= 0 {
		query.Offset(int64(offset))
	}
	if limit > 0 {
		query.Limit(int64(limit))
	}
	if order != "" {
		query.OrderBy(order)
	}

	err := query.All(&result)
	return result, err
}
