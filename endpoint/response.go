package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// nolint:godox TODO: move to smis

//Response provides functions to write http responses
type Response struct {
	Log logrus.FieldLogger
}

func (r *Response) logError(msg string) {
	if r == nil || r.Log == nil {
		return
	}

	r.Log.Error(msg)
}

//WriteJSON sends a JSON response with given code and payload
func (r *Response) WriteJSON(writer http.ResponseWriter, code int, payload interface{}) {
	if r == nil {
		return
	}

	response, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Sprintf("failed to encode response payload: %v", err)
		r.logError(msg)
		writer.WriteHeader(http.StatusInternalServerError) // nolint:godox TODO: find a way to write Default Error, maybe header?

		if _, err := writer.Write([]byte(msg)); err != nil {
			r.logError(fmt.Sprintf("failed to write response: %v", err))
		}

		return
	}

	writer.WriteHeader(code)

	if _, err := writer.Write(response); err != nil {
		r.logError(fmt.Sprintf("failed to write response: %v", err))
	}
}
