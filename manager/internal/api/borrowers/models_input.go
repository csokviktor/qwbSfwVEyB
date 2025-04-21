package borrowers

import (
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"github.com/gin-gonic/gin"
)

type BorrowerCreateValidator struct {
	Name string `form:"name" json:"name" binding:"required,min=1"`
}

func (uc *BorrowerCreateValidator) Bind(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(uc)
}

func (uc *BorrowerCreateValidator) GetBorrower() *dbmodels.Borrower {
	return &dbmodels.Borrower{Name: uc.Name}
}
