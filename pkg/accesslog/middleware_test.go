package accesslog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/pkg/log"
)

func TestHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://127.0.0.1/v1/deposits/balance", strings.NewReader(`{"owner_id":"11111111-1111-1111-1111-111111111111"`))
	ctx := routing.NewContext(res, req)

	logger, entries := log.NewForTest()
	handler := Handler(logger)
	err := handler(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, entries.Len())
	assert.Equal(t, "POST /v1/deposits/balance HTTP/1.1 200 0", entries.All()[0].Message)
}
