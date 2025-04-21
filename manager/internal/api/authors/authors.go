package authors

import (
	"errors"
	"net/http"

	"github.com/csokviktor/lib_manager/internal/api"
	"github.com/csokviktor/lib_manager/internal/repository"
	"github.com/csokviktor/lib_manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type routes struct {
	authorsService service.Authors
}

type Routes interface {
	CreateAuthor(ctx *gin.Context)
	GetAuthor(ctx *gin.Context)
	GetAuthors(ctx *gin.Context)
}

func NewRoutes(authorsService service.Authors) Routes {
	return &routes{
		authorsService,
	}
}

func (r *routes) CreateAuthor(ctx *gin.Context) {
	validator := AuthorCreateValidator{}
	if err := validator.Bind(ctx); err != nil {
		log.Error().Msg("CreateAuthor validation failed")
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("CreateAuthor")
	newAuthor, err := r.authorsService.Create(ctx, validator.GetAuthor())
	if err != nil {
		log.Error().Msg("CreateAuthor failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, DBToAPI(newAuthor))
}

func (r *routes) GetAuthor(ctx *gin.Context) {
	authorID := ctx.Param("id")
	if err := uuid.Validate(authorID); err != nil {
		log.Error().Msgf("GetAuthor id %s validation failed", authorID)
		ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err))
		return
	}

	log.Debug().Msgf("GetAuthor with %s id", authorID)
	author, err := r.authorsService.GetByID(ctx, authorID)
	if err != nil {
		log.Error().Msgf("GetAuthor failed with id %s", authorID)
		if errors.Is(err, repository.NotFoundError{}) {
			ctx.JSON(http.StatusNotFound, api.NewErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBToAPI(author))
}

func (r *routes) GetAuthors(ctx *gin.Context) {
	log.Debug().Msgf("GetAuthors")
	authors, err := r.authorsService.GetAll(ctx)
	if err != nil {
		log.Error().Msg("GetAuthors failed")
		ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, DBsToAPIs(authors))
}
