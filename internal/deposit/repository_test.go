package deposit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/test"
	"users-balance-microservice/pkg/log"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "deposit")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	ownerId := uuid.New()
	dep := entity.Deposit{OwnerId: ownerId, Balance: 1000}

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create deposit
	err = repo.Create(ctx, dep)
	if assert.NoError(t, err) {
		count2, _ := repo.Count(ctx)
		assert.EqualValues(t, 1, count2-count)
	}

	// get balance
	dep, err = repo.Get(ctx, ownerId)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1000, dep.Balance)
	}

	// update balance
	dep.Balance -= 600
	err = repo.Update(ctx, dep)
	if assert.NoError(t, err) {
		dep, _ = repo.Get(ctx, ownerId)
		assert.EqualValues(t, 400, dep.Balance)
	}

	// push an update with negative balance -> get an error, update rejected
	dep.Balance -= 20000
	err = repo.Update(ctx, dep)
	if assert.Error(t, err) {
		dep, _ = repo.Get(ctx, ownerId)
		assert.EqualValues(t, 400, dep.Balance)
	}

}
