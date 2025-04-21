package repository

import (
	"context"
	"errors"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"gorm.io/gorm"
)

type borrowers struct {
	db *gorm.DB
}

//go:generate mockgen -source borrowers.go -destination mocks/borrowers.go --package mocks

type Borrowers interface {
	CRUD[dbmodels.Borrower]
}

func NewBorrowers(db *gorm.DB) Borrowers {
	return &borrowers{
		db,
	}
}

func (b *borrowers) Create(ctx context.Context, borrower *dbmodels.Borrower) (*dbmodels.Borrower, error) {
	result := b.db.WithContext(ctx).Create(borrower)
	return borrower, result.Error
}

func (b *borrowers) GetAll(ctx context.Context) ([]dbmodels.Borrower, error) {
	var borrowers []dbmodels.Borrower
	result := b.db.WithContext(ctx).Preload(bookClause).Find(&borrowers)
	return borrowers, result.Error
}

func (b *borrowers) GetByID(ctx context.Context, id string) (*dbmodels.Borrower, error) {
	var borrower dbmodels.Borrower
	result := b.db.WithContext(ctx).Preload(bookClause).First(&borrower, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFoundByID(id)
	}
	return &borrower, result.Error
}

func (b *borrowers) Update(ctx context.Context, borrower *dbmodels.Borrower) (*dbmodels.Borrower, error) {
	if _, err := b.GetByID(ctx, borrower.ID); err != nil {
		return nil, err
	}
	result := b.db.WithContext(ctx).Updates(borrower)
	return borrower, result.Error
}

func (b *borrowers) Delete(ctx context.Context, id string) error {
	return b.db.WithContext(ctx).Delete(&dbmodels.Borrower{}, id).Error
}
