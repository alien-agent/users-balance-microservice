package deposit

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/deposits/balance", res.getBalance)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getBalance(c *routing.Context) error {
	var input GetBalanceRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	deposit, err := r.service.Get(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.Write(deposit.Balance)
}
