package service

import (
	"context"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
)

type authors struct {
	authorsRepository repository.Authors
}

//go:generate mockgen -source authors.go -destination mocks/authors.go --package mocks

type Authors interface {
	Create(ctx context.Context, newAuthor *dbmodels.Author) (*dbmodels.Author, error)
	GetByID(ctx context.Context, id string) (*dbmodels.Author, error)
	GetAll(ctx context.Context) ([]dbmodels.Author, error)
}

func NewAuthors(authorsRepository repository.Authors) Authors {
	return &authors{
		authorsRepository,
	}
}

func (a *authors) Create(ctx context.Context, newAuthor *dbmodels.Author) (*dbmodels.Author, error) {
	createdAuthor, err := a.authorsRepository.Create(ctx, newAuthor)
	return createdAuthor, err
}

func (a *authors) GetByID(ctx context.Context, id string) (*dbmodels.Author, error) {
	author, err := a.authorsRepository.GetByID(ctx, id)
	return author, err
}

func (a *authors) GetAll(ctx context.Context) ([]dbmodels.Author, error) {
	allAuthors, err := a.authorsRepository.GetAll(ctx)
	return allAuthors, err
}
