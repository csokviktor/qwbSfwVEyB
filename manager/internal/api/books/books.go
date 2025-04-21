package books

import (
	"errors"
	"net/http"

	"github.com/csokviktor/qwbSfwVEyB/manager/internal/api"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type routes struct {
	booksService service.Books
}

type Routes interface {
	CreateBook(ctx *gin.Context)
	GetBooks(ctx *gin.Context)
	BorrowBook(ctx *gin.Context)
}

func NewRoutes(booksService service.Books) Routes {
	return &routes{
		booksService,
	}
}

func (r *routes) CreateBook(ctx *gin.Context) {
	validator := BookCreateValidator{}
	if err := validator.Bind(ctx); err != nil {
		log.Error().Msg("CreateBook validation failed")
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("CreateBook")
	newBook, err := r.booksService.Create(ctx, validator.GetBook())
	if err != nil {
		log.Error().Msg("CreateBook failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, DBToAPI(newBook))
}

func (r *routes) GetBooks(ctx *gin.Context) {
	log.Debug().Msgf("GetBooks")
	books, err := r.booksService.GetAll(ctx)
	if err != nil {
		log.Error().Msg("GetBooks failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBsToAPIs(books))
}

func (r *routes) BorrowBook(ctx *gin.Context) {
	validator := BookBorrowValidator{}
	if err := validator.Bind(ctx); err != nil {
		log.Error().Msg("BorrowBook validation failed")
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	bookID := ctx.Param("id")
	if err := uuid.Validate(bookID); err != nil {
		log.Error().Msgf("BorrowBook id %s validation failed", bookID)
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("BorrowBook with %s id for borrower with id %s", bookID, validator.BorrowerID)
	if err := r.booksService.Borrow(ctx, validator.BorrowerID, bookID); err != nil {
		log.Error().Msgf("BorrowBook failed with %s id for borrower with id %s", bookID, validator.BorrowerID)
		if errors.Is(err, repository.NotFoundError{}) {
			ctx.JSON(http.StatusNotFound, api.NewErrorResponse(err))
			return
		} else if errors.Is(err, service.WrongArgumentError{}) {
			ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}
