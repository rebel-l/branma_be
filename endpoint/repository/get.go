package repository

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rebel-l/branma_be/endpoint"
)

// Get returns a repository identified by ID
func (h *Handler) Get(writer http.ResponseWriter, request *http.Request) {
	response := endpoint.Response{}
	payload := &Payload{}

	// 0. validate request
	if request == nil {
		payload.Error = "request is empty"
		response.WriteJSON(writer, http.StatusBadRequest, payload)

		return
	}

	response.Log = h.svc.NewLogForRequestID(request.Context())

	idRaw, ok := mux.Vars(request)["id"]
	if !ok {
		payload.Error = "id must be given"
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
	if err != nil {
		payload.Error = fmt.Sprintf("failed to load repository for id: %d", id)
		response.WriteJSON(writer, http.StatusInternalServerError, payload)

		return
	}

	// 2. send response
	payload.Repository = model
	response.WriteJSON(writer, http.StatusOK, payload)
}
