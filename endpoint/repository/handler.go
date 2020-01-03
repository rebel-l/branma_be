package repository

import (
	"fmt"
	"net/http"

	"github.com/rebel-l/branma_be/repository/repositorymapper"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/smis"
)

const (
	errRequestEmpty = "request is empty"
	errNoID         = "id must be given"
)

// Handler provides useful variables for the specific endpoint handlers
type Handler struct {
	svc    *smis.Service
	mapper *repositorymapper.Mapper // nolint:godox TODO: change to interface
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

	_, err := svc.RegisterEndpoint("/repository/{id}", http.MethodGet, endpoint.get)
	if err != nil {
		return fmt.Errorf("failed to init get endpoint for repository: %w", err)
	}

	_, err = svc.RegisterEndpoint("/repository", http.MethodPut, endpoint.put)
	if err != nil {
		return fmt.Errorf("failed to init put endpoint for repository: %w", err)
	}

	_, err = svc.RegisterEndpoint("/repository/{id}", http.MethodDelete, endpoint.delete)
	if err != nil {
		return fmt.Errorf("failed to init delete endpoint for repository: %w", err)
	}

	return err
}
