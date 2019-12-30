package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rebel-l/branma_be/repository/repositorymapper"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/smis"
)

type Handler struct {
	svc    *smis.Service
	mapper *repositorymapper.Mapper
}

func New(svc *smis.Service, db *sqlx.DB) *Handler {
	return &Handler{
		svc:    svc,
		mapper: repositorymapper.New(db),
	}
}

func Init(svc *smis.Service, db *sqlx.DB) error {
	endpoint := New(svc, db)

	_, err := svc.RegisterEndpoint("/repository/{id}", http.MethodGet, endpoint.get)
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

func (h *Handler) get(writer http.ResponseWriter, request *http.Request) {
	log := h.svc.NewLogForRequestID(request.Context())

	idRaw, ok := mux.Vars(request)["id"]
	if !ok {
		msg := "id must be given"
		log.Errorf(msg)
		writer.WriteHeader(http.StatusBadRequest)
		if _, err := writer.Write([]byte(msg)); err != nil {
			log.Errorf("failed to write response: %v", err)
		}
	}

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		msg := "converting id to integer failed"
		log.Errorf("%s: %s", msg, err)
		writer.WriteHeader(http.StatusBadRequest)
		if _, err := writer.Write([]byte(msg)); err != nil {
			log.Errorf("failed to write response: %v", err)
		}
	}

	repo, err := h.mapper.Load(request.Context(), id)
	if err != nil {
		msg := fmt.Sprintf("failed to load repository for id: %d", id)
		log.Errorf("%s: %s", msg, err)
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(msg)); err != nil {
			log.Errorf("failed to write response: %v", err)
		}
	}

	j, err := json.Marshal(repo)
	if err != nil {
		msg := fmt.Sprintf("failed to convert repository for id: %d to JSON", id)
		log.Errorf("%s: %s", msg, err)
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write([]byte(msg)); err != nil {
			log.Errorf("failed to write response: %v", err)
		}
	}

	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(j); err != nil {
		log.Errorf("failed to write response: %v", err)
	}
}

func (h *Handler) delete(writer http.ResponseWriter, request *http.Request) {
	log := h.svc.NewLogForRequestID(request.Context())
	log.Info("endpoint not implemented yet")
}
