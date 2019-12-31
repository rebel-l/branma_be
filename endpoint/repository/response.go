package repository

import "github.com/rebel-l/branma_be/repository/repositorymodel"

// Payload represents response payload for endpoint
type Payload struct {
	Repository *repositorymodel.Repository `json:"repository,omitempty"`
	Error      string                      `json:"error,omitempty"`
}

// NewPayload returns a new Payload struct
func NewPayload(repository *repositorymodel.Repository) *Payload {
	return &Payload{Repository: repository}
}
