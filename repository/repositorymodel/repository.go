package repositorymodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	// ErrDecodeJSON occurs if the a string is not in JSON format
	ErrDecodeJSON = errors.New("failed to decode JSON")
)

// Repository represents a model of repository including business logic
type Repository struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// DecodeJSON converts JSON data to struct
func (r *Repository) DecodeJSON(reader io.Reader) error {
	if r == nil {
		return nil
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(r); err != nil {
		return fmt.Errorf("%w: %v", ErrDecodeJSON, err)
	}

	return nil
}
