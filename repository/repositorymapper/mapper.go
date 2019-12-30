package repositorymapper

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/branma_be/repository/repositorymodel"
	"github.com/rebel-l/branma_be/repository/repositorystore"
)

var (
	// ErrLoadFromDB occurs if something went wrong on loading
	ErrLoadFromDB = errors.New("failed to load repository from database")

	// ErrNoData occurs if given model is nil
	ErrNoData = errors.New("repository is nil")

	// ErrSaveToDB occurs if something went wrong on saving
	ErrSaveToDB = errors.New("failed to save repository to database")

	// ErrDeleteFromDB occurs if something went wrong on deleting
	ErrDeleteFromDB = errors.New("failed to delete repository from database")
)

// Mapper provides methods to load and persist repository models
type Mapper struct {
	db *sqlx.DB
}

// New returns a new mapper
func New(db *sqlx.DB) *Mapper {
	return &Mapper{db: db}
}

// Load returns a repository model loaded from database by ID
func (m *Mapper) Load(ctx context.Context, id int) (*repositorymodel.Repository, error) {
	s := &repositorystore.Repository{ID: id}
	if err := s.Read(ctx, m.db); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadFromDB, err)
	}

	return storeToModel(s), nil
}

// Save persists (create or update) the model and returns the changed data (id, createdAt or modifiedAt)
func (m *Mapper) Save(ctx context.Context, model *repositorymodel.Repository) (*repositorymodel.Repository, error) {
	if model == nil {
		return nil, ErrNoData
	}

	s := modelToStore(model)

	if model.ID != 0 {
		if err := s.Update(ctx, m.db); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSaveToDB, err)
		}
	} else {
		if err := s.Create(ctx, m.db); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSaveToDB, err)
		}
	}

	model = storeToModel(s)

	return model, nil
}

// Delete removes a model from database by ID
func (m *Mapper) Delete(ctx context.Context, id int) error {
	s := &repositorystore.Repository{ID: id}
	if err := s.Delete(ctx, m.db); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFromDB, err)
	}

	return nil
}

func storeToModel(s *repositorystore.Repository) *repositorymodel.Repository {
	if s == nil {
		return &repositorymodel.Repository{}
	}

	return &repositorymodel.Repository{
		ID:         s.ID,
		Name:       s.Name,
		URL:        s.URL,
		CreatedAt:  s.CreatedAt,
		ModifiedAt: s.ModifiedAt,
	}
}

func modelToStore(m *repositorymodel.Repository) *repositorystore.Repository {
	if m == nil {
		return &repositorystore.Repository{}
	}

	return &repositorystore.Repository{
		ID:         m.ID,
		Name:       m.Name,
		URL:        m.URL,
		CreatedAt:  m.CreatedAt,
		ModifiedAt: m.ModifiedAt,
	}
}
