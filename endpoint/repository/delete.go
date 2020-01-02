package repository

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rebel-l/smis"
)

// Delete removes a repository identified by ID
func (h *Handler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	// 1. delete model
	if err := h.mapper.Delete(request.Context(), id); err != nil {
		response.Log.Error(err)

		payload.Error = fmt.Sprintf("failed to delete repository for id: %d", id)
		response.WriteJSON(writer, http.StatusInternalServerError, payload)

		return
	}

	// 2. send response
	response.WriteJSON(writer, http.StatusOK, payload)
}
