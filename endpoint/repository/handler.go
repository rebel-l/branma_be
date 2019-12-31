package repository

import (
	"fmt"
	"net/http"

	"github.com/rebel-l/branma_be/repository/repositorymapper"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/smis"
)

// Handler provides useful variables for the specific endpoint handlers
type Handler struct {
	svc    *smis.Service
	mapper *repositorymapper.Mapper
}

// New returns a new handler
func New(svc *smis.Service, db *sqlx.DB) *Handler {
	return &Handler{
		svc:    svc,
		mapper: repositorymapper.New(db),
	}
}

// Init initialises the endpoints for the repository
func Init(svc *smis.Service, db *sqlx.DB) error {
	endpoint := New(svc, db)

	_, err := svc.RegisterEndpoint("/repository/{id}", http.MethodGet, endpoint.Get)
	if err != nil {
		return fmt.Errorf("failed to init get endpoint for repository: %w", err)
	}

	_, err = svc.RegisterEndpoint("/repository", http.MethodPut, endpoint.Put)
	if err != nil {
		return fmt.Errorf("failed to init put endpoint for repository: %w", err)
	}

	_, err = svc.RegisterEndpoint("/repository/{id}", http.MethodDelete, endpoint.delete)
	if err != nil {
		return fmt.Errorf("failed to init delete endpoint for repository: %w", err)
	}

	return err
}

func (h *Handler) delete(writer http.ResponseWriter, request *http.Request) {
	log := h.svc.NewLogForRequestID(request.Context())
	log.Info("endpoint not implemented yet")
}
