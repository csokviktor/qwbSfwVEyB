package borrowers

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
	borrowersService service.Borrowers
}

type Routes interface {
	CreateBorrower(ctx *gin.Context)
	GetBorrower(ctx *gin.Context)
	GetBorrowers(ctx *gin.Context)
	GetBorrowerBooks(ctx *gin.Context)
}

func NewRoutes(borrowersService service.Borrowers) Routes {
	return &routes{
		borrowersService,
	}
}

func (r *routes) CreateBorrower(ctx *gin.Context) {
	validator := BorrowerCreateValidator{}
	if err := validator.Bind(ctx); err != nil {
		log.Error().Msg("CreateBorrower validation failed")
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("CreateBorrower")
	newBorrower, err := r.borrowersService.Create(ctx, validator.GetBorrower())
	if err != nil {
		log.Error().Msg("CreateBorrower failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, DBToAPI(newBorrower))
}

func (r *routes) GetBorrower(ctx *gin.Context) {
	borrowerID := ctx.Param("id")
	if err := uuid.Validate(borrowerID); err != nil {
		log.Error().Msgf("GetBorrower id %s validation failed", borrowerID)
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("GetBorrower with %s id", borrowerID)
	borrower, err := r.borrowersService.GetByID(ctx, borrowerID)
	if err != nil {
		log.Error().Msgf("GetBorrower failed with id %s", borrowerID)
		if errors.Is(err, repository.NotFoundError{}) {
			ctx.JSON(http.StatusNotFound, api.NewErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBToAPI(borrower))
}

func (r *routes) GetBorrowers(ctx *gin.Context) {
	log.Debug().Msgf("GetBorrowers")
	borrowers, err := r.borrowersService.GetAll(ctx)
	if err != nil {
		log.Error().Msg("GetBorrowers failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBsToAPIs(borrowers))
}

func (r *routes) GetBorrowerBooks(ctx *gin.Context) {
	borrowerID := ctx.Param("id")
	if err := uuid.Validate(borrowerID); err != nil {
		log.Error().Msgf("GetBorrowerBooks id %s validation failed", borrowerID)
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("GetBorrowerBooks with %s id", borrowerID)
	borrower, err := r.borrowersService.GetByID(ctx, borrowerID)
	if err != nil {
		log.Error().Msgf("GetBorrowerBooks failed with id %s", borrowerID)
		if errors.Is(err, repository.NotFoundError{}) {
			ctx.JSON(http.StatusNotFound, api.NewErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBToAPI(borrower).Books)
}
