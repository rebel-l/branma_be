package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rebel-l/branma_be/endpoint"
	"github.com/rebel-l/branma_be/repository/repositorymodel"
)

func (h *Handler) Put(writer http.ResponseWriter, request *http.Request) {
	response := endpoint.Response{}
	payload := &Payload{}
	// 0. validate request
	if request == nil {
		payload.Error = "request is empty"
		response.WriteJSON(writer, http.StatusBadRequest, payload)

		return
	}

	response.Log = h.svc.NewLogForRequestID(request.Context())

	if request.Body == nil {
		payload.Error = fmt.Sprint("request body is empty")
		response.WriteJSON(writer, http.StatusBadRequest, payload)
		return
	}

	// 1. decode payload
	model := &repositorymodel.Repository{}
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(model); err != nil {
		payload.Error = fmt.Sprintf("failed to decode request payload: %v", err)
		response.WriteJSON(writer, http.StatusInternalServerError, payload)

		return
	}

	code := http.StatusOK
	if model.ID == 0 {
		code = http.StatusCreated
	}

	// 2. save model
	model, err := h.mapper.Save(request.Context(), model)
	if err != nil {
		payload.Error = fmt.Sprintf("failed to save repository: %v", err)
		response.WriteJSON(writer, http.StatusInternalServerError, payload)

		return
	}

	// 3. send response
	payload.Repository = model
	response.WriteJSON(writer, code, payload)
}
