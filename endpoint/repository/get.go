package repository

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rebel-l/branma_be/repository/repositorymapper"
	"github.com/rebel-l/smis"

	"github.com/gorilla/mux"
)

// get returns a repository identified by ID
func (h *Handler) get(writer http.ResponseWriter, request *http.Request) {
	response := smis.Response{}
	payload := &Payload{}

	// 0. validate request
	if request == nil {
		payload.Error = errRequestEmpty
		response.WriteJSON(writer, http.StatusBadRequest, payload)

		return
	}

	response.Log = h.svc.NewLogForRequestID(request.Context())

	idRaw, ok := mux.Vars(request)["id"]
	if !ok {
		payload.Error = errNoID
		response.WriteJSON(writer, http.StatusBadRequest, payload)

		return
	}

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		payload.Error = "converting id to integer failed"
		response.WriteJSON(writer, http.StatusBadRequest, payload)

		return
	}

	// 1. load model
	model, err := h.mapper.Load(request.Context(), id)
	if errors.Is(err, repositorymapper.ErrNotFound) {
		payload.Error = fmt.Sprintf("repository with id %d not found", id)
		response.WriteJSON(writer, http.StatusNotFound, payload)

		return
	} else if err != nil {
		response.Log.Error(err)

		payload.Error = fmt.Sprintf("failed to load repository for id: %d", id)
		response.WriteJSON(writer, http.StatusInternalServerError, payload)

		return
	}

	// 2. send response
	payload.Repository = model
	response.WriteJSON(writer, http.StatusOK, payload)
}
