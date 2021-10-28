package deposit

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/test"
	"users-balance-microservice/pkg/log"
)

const invalidIdResponse = `{"status":400,"message":"There is some problem with the data you submitted.","details":[{"field":"owner_id","error":"must be a valid UUID"}]}`
const badRequestResponse = `{"status":400,"message":"Your request is in a bad format."}`

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{
		items: []entity.Deposit{
			{uuid.MustParse("615f3e76-37d3-11ec-8d3d-0242ac130003"), 1000},
		},
	}
	RegisterHandlers(router.Group(""), NewService(repo, logger), logger)

	tests := []test.APITestCase{
		{"get balance existing", "POST", "/deposits/balance", `{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003"}`, http.StatusOK, `1000`},
		{"get balance non-existing", "POST", "/deposits/balance", `{"owner_id": "8c5593a0-37d3-11ec-8d3d-0242ac130003"}`, http.StatusOK, `0`},
		{"get balance invalid owner_id", "POST", "/deposits/balance", `{"owner_id": "0123456789"}`, http.StatusBadRequest, invalidIdResponse},
		{"get balance invalid request", "POST", "/deposits/balance", `{"owner_id": `, http.StatusBadRequest, badRequestResponse},
		{"get balance invalid method", "GET", "/deposits/balance", "", http.StatusMethodNotAllowed, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
