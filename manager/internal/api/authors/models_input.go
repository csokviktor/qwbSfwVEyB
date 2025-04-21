package authors

import (
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"github.com/gin-gonic/gin"
)

type AuthorCreateValidator struct {
	Name string `form:"name" json:"name" binding:"required,min=1"`
}

func (uc *AuthorCreateValidator) Bind(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(uc)
}

func (uc *AuthorCreateValidator) GetAuthor() *dbmodels.Author {
	return &dbmodels.Author{Name: uc.Name}
}
