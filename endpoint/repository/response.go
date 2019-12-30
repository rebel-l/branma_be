package repository

import "github.com/rebel-l/branma_be/repository/repositorymodel"

type Payload struct {
	Repository *repositorymodel.Repository `json:"repository,omitempty"`
	Error      string                      `json:"error,omitempty"`
}

func NewPayload(repository *repositorymodel.Repository) *Payload {
	return &Payload{Repository: repository}
}

// TODO: from JSON
// TODO: to JSON
