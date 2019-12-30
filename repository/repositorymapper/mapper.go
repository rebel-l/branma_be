package repositorymapper

import (
	"context"
	"errors"
	"fmt"

	"github.com/rebel-l/branma_be/repository/repositorymodel"

	"github.com/rebel-l/branma_be/repository/repositorystore"

	"github.com/jmoiron/sqlx"
)

var (
	ErrLoadFromDB error = errors.New("failed to load repository from database")
	ErrNoData     error = errors.New("repository is nil")
	ErrSaveToDB   error = errors.New("failed to save repository to database")
)

type Mapper struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Mapper {
	return &Mapper{db: db}
}

func (m *Mapper) Load(ctx context.Context, id int) (*repositorymodel.Repository, error) {
	s := &repositorystore.Repository{ID: id}
	if err := s.Read(ctx, m.db); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadFromDB, err)
	}

	return storeToModel(s), nil
}

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
