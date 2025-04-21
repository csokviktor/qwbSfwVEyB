package service

import (
	"context"

	"github.com/csokviktor/lib_manager/internal/repository"
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
)

type borrowers struct {
	borrowersRepository repository.Borrowers
}

//go:generate mockgen -source borrowers.go -destination mocks/borrowers.go --package mocks

type Borrowers interface {
	Create(ctx context.Context, newBorrower *dbmodels.Borrower) (*dbmodels.Borrower, error)
	GetByID(ctx context.Context, id string) (*dbmodels.Borrower, error)
	GetAll(ctx context.Context) ([]dbmodels.Borrower, error)
}

func NewBorrowers(borrowersRepository repository.Borrowers) Borrowers {
	return &borrowers{
		borrowersRepository,
	}
}

func (b *borrowers) Create(ctx context.Context, newBorrower *dbmodels.Borrower) (*dbmodels.Borrower, error) {
	createdBorrower, err := b.borrowersRepository.Create(ctx, newBorrower)
	return createdBorrower, err
}

func (b *borrowers) GetByID(ctx context.Context, id string) (*dbmodels.Borrower, error) {
	borrower, err := b.borrowersRepository.GetByID(ctx, id)
	return borrower, err
}

func (b *borrowers) GetAll(ctx context.Context) ([]dbmodels.Borrower, error) {
	allBorrowers, err := b.borrowersRepository.GetAll(ctx)
	return allBorrowers, err
}
