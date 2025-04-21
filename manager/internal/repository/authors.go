package repository

import (
	"context"
	"errors"
	"time"

	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"github.com/sony/gobreaker"
	"gorm.io/gorm"
)

type authors struct {
	db      *gorm.DB
	breaker *gobreaker.CircuitBreaker
}

//go:generate mockgen -source authors.go -destination mocks/authors.go --package mocks

type Authors interface {
	CRUD[dbmodels.Author]
}

func NewAuthors(db *gorm.DB) Authors {
	cbSettings := gobreaker.Settings{
		Name:        "AuthorsRepositoryBreaker",
		MaxRequests: 5,
		Interval:    0,               // no rolling window reset
		Timeout:     5 * time.Second, // how long the circuit stays open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	return &authors{
		db:      db,
		breaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func (a *authors) Create(ctx context.Context, author *dbmodels.Author) (*dbmodels.Author, error) {
	result := a.db.WithContext(ctx).Create(author)
	return author, result.Error
}

func (a *authors) GetAll(ctx context.Context) ([]dbmodels.Author, error) {
	// circuit breaker demo for DB.
	// in sqlite this does not make sense so I will not implement it everywhere
	// but could be useful with remote dbs
	result, err := a.breaker.Execute(func() (interface{}, error) {
		var authors []dbmodels.Author
		dbResult := a.db.WithContext(ctx).Preload(bookClause).Find(&authors)
		return authors, dbResult.Error
	})
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // value is always []dbmodels.Author if exists
	return result.([]dbmodels.Author), nil
}

func (a *authors) GetByID(ctx context.Context, id string) (*dbmodels.Author, error) {
	var author dbmodels.Author
	result := a.db.WithContext(ctx).Preload(bookClause).First(&author, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFoundByID(id)
	}
	return &author, result.Error
}

func (a *authors) Update(ctx context.Context, author *dbmodels.Author) (*dbmodels.Author, error) {
	if _, err := a.GetByID(ctx, author.ID); err != nil {
		return nil, err
	}
	result := a.db.WithContext(ctx).Updates(author)
	return author, result.Error
}

func (a *authors) Delete(ctx context.Context, id string) error {
	return a.db.WithContext(ctx).Delete(&dbmodels.Author{}, id).Error
}
