package service

import (
	"context"
	"fmt"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"github.com/rs/zerolog/log"
)

type books struct {
	booksRepository  repository.Books
	authorsService   Authors
	borrowersService Borrowers
}

//go:generate mockgen -source books.go -destination mocks/books.go --package mocks

type Books interface {
	Create(ctx context.Context, newBook *dbmodels.Book) (*dbmodels.Book, error)
	GetAll(ctx context.Context) ([]dbmodels.Book, error)
	Borrow(ctx context.Context, borrowerID, bookID string) error
}

func NewBooks(
	bookRepository repository.Books,
	authorsService Authors,
	borrowersService Borrowers,
) Books {
	return &books{
		bookRepository,
		authorsService,
		borrowersService,
	}
}

func (b *books) Create(ctx context.Context, newBook *dbmodels.Book) (*dbmodels.Book, error) {
	log.Debug().Msgf("GetByID for author with id %s", newBook.AuthorID)
	_, err := b.authorsService.GetByID(ctx, newBook.AuthorID)
	if err != nil {
		return nil, err
	}
	createdBook, err := b.booksRepository.Create(ctx, newBook)
	return createdBook, err
}

func (b *books) GetAll(ctx context.Context) ([]dbmodels.Book, error) {
	books, err := b.booksRepository.GetAll(ctx)
	return books, err
}

func (b *books) Borrow(ctx context.Context, borrowerID, bookID string) error {
	log.Debug().Msgf("GetByID for borrower with id %s", borrowerID)
	borrower, err := b.borrowersService.GetByID(ctx, borrowerID)
	if err != nil {
		return err
	}

	log.Debug().Msgf("GetByID for book with id %s", borrowerID)
	book, err := b.booksRepository.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book.BorrowerID != nil {
		log.Warn().Msgf("book with %s id already borrowed", book.ID)
		return fmt.Errorf("book with %s id already borrowed, %w", book.ID, WrongArgumentError{})
	}

	book.BorrowerID = &borrower.ID
	_, err = b.booksRepository.Update(ctx, book)
	return err
}
