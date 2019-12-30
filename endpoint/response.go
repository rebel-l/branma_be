package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// TODO: move to smis

type Response struct {
	Log logrus.FieldLogger
}

func (r *Response) logError(msg string) {
	if r == nil || r.Log == nil {
		return
	}

	r.Log.Error(msg)
}

func (r *Response) WriteJSON(writer http.ResponseWriter, code int, payload interface{}) {
	if r == nil {
		return
	}
	response, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Sprintf("failed to encode response payload: %v", err)
		r.logError(msg)
		writer.WriteHeader(http.StatusInternalServerError) // TODO: find a way to write Default Error, maybe header?
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
