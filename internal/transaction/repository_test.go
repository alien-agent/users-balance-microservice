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

	// create (deposit withdrawal)
	tx := entity.Transaction{
		Id:              0,
		SenderId:        id1,
		RecipientId:     uuid.Nil,
		Amount:          500,
		Description:     "Monthly subscription",
		TransactionDate: time.Now(),
	}
	err := repo.Create(ctx, &tx)
	assert.Nil(t, err)

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
	assert.Nil(t, err)
	assert.NotEqual(t, 0, tx.Id) // tx.Id should be auto-incremented by database

	// create (money-transfer)
	tx = entity.Transaction{
		Id:              0,
		SenderId:        id1,
		RecipientId:     id2,
		Amount:          1500,
		Description:     "thanks for dinner!",
		TransactionDate: time.Now(),
	}
	err = repo.Create(ctx, &tx)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, tx.Id) // tx.Id should be auto-incremented by database

	// create with negative amount -> db error
	err = repo.Create(ctx, &entity.Transaction{
		Id:              0,
		SenderId:        id2,
		RecipientId:     id1,
		Amount:          -1000,
		Description:     "happy birthday!",
		TransactionDate: time.Now(),
	})
	assert.NotNil(t, err)
}
