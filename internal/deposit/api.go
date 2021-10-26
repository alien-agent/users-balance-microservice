package deposit

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/google/uuid"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/deposits/<owner_id>", res.get)
	r.Patch("/deposits/<owner_id>", res.update)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	var ownerId uuid.UUID
	ownerId, err := uuid.Parse(c.Param("owner_id"))
	if err != nil {
		return errors.BadRequest("No valid UUID received.")
	}

	deposit, err := r.service.Get(c.Request.Context(), ownerId)
	if err != nil {
		return err
	}

	return c.Write(deposit.Balance)
}

func (r resource) update(c *routing.Context) error {
	var ownerId uuid.UUID
	ownerId, err := uuid.Parse(c.Param("owner_id"))
	if err != nil {
		return errors.BadRequest("No valid UUID received.")
	}

	deposit, err := r.service.Get(c.Request.Context(), ownerId)
	if err != nil {
		return err
	}

	return c.Write(deposit.Balance)
}
