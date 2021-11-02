package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/test"
	"users-balance-microservice/pkg/log"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "transaction")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	id1, id2 := uuid.New(), uuid.New()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create (deposit withdrawal)
	tx := entity.Transaction{
		Id:              0,
		SenderId:        id1,
		RecipientId:     uuid.Nil,
		Amount:          300,
		Description:     "Monthly subscription",
		TransactionDate: time.Now(),
	}
	err = repo.Create(ctx, &tx)
	if assert.NoError(t, err) {
		count2, err := repo.Count(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), count2-count) // new transaction is really in db
			count++
		}
	}

	// create (deposit top-up)
	tx = entity.Transaction{
		Id:              0,
		SenderId:        uuid.Nil,
		RecipientId:     id1,
		Amount:          500,
		Description:     "VISA top-up",
		TransactionDate: time.Now(),
	}
	err = repo.Create(ctx, &tx)
	if assert.NoError(t, err) {
		assert.NotZero(t, tx.Id) // tx.Id should be auto-incremented by db.

		count2, err := repo.Count(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), count2-count) // new transaction is really in db
			count++
		}
	}

	// create (money transfer)
	tx = entity.Transaction{
		Id:              0,
		SenderId:        id1,
		RecipientId:     id2,
		Amount:          1500,
		Description:     "thanks for dinner!",
		TransactionDate: time.Now(),
	}
	err = repo.Create(ctx, &tx)
	if assert.NoError(t, err) {
		assert.NotZero(t, tx.Id)

		count2, err := repo.Count(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(1), count2-count) // new transaction is really in db
			count++
		}
	}

	// create with negative amount -> db error
	err = repo.Create(ctx, &entity.Transaction{
		Id:              0,
		SenderId:        id2,
		RecipientId:     id1,
		Amount:          -1000,
		Description:     "happy birthday!",
		TransactionDate: time.Now(),
	})
	if assert.Error(t, err) {
		count2, err := repo.Count(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, int64(0), count2-count) // invalid transaction is NOT in db
		}
	}

	// list for user
	txs, err := repo.GetForUser(ctx, id1, "", 0, 0)
	if assert.NoError(t, err) {
		assert.Len(t, txs, 3)
	}

	// list for user with pagination
	txs, err = repo.GetForUser(ctx, id1, "", 1, 1)
	if assert.NoError(t, err) {
		assert.Len(t, txs, 1)
	}
}
