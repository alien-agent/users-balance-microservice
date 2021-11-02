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
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get balance
	dep, err = repo.Get(ctx, ownerId)
	assert.Nil(t, err)
	assert.Equal(t, int64(1000), dep.Balance)
	/*_, err = depositRepo.GetBalance(ctx, uuid.MustParse("11f58ca1-8fee-453a-8bf0-544b4bcde3f2"))
	assert.Equal(t, sql.ErrNoRows, err)*/

	// update balance
	dep.Balance -= 600
	err = repo.Update(ctx, dep)
	assert.Nil(t, err)
	dep, _ = repo.Get(ctx, ownerId)
	assert.Equal(t, int64(400), dep.Balance)

	// push an update with negative balance -> get an error, update rejected
	dep.Balance -= 20000
	err = repo.Update(ctx, dep)
	assert.NotNil(t, err)
	dep, _ = repo.Get(ctx, ownerId)
	assert.Equal(t, int64(400), dep.Balance)
}
