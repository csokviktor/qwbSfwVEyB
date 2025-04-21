package books

import (
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"github.com/gin-gonic/gin"
)

type BookCreateValidator struct {
	Title    string `form:"title" json:"title" binding:"required,min=1"`
	AuthorID string `form:"authorID" json:"authorID" binding:"required,uuid"`
}

func (bc *BookCreateValidator) Bind(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(bc)
}

func (bc *BookCreateValidator) GetBook() *dbmodels.Book {
	return &dbmodels.Book{Title: bc.Title, AuthorID: bc.AuthorID}
}

type BookBorrowValidator struct {
	BorrowerID string `form:"borrowerID" json:"borrowerID" binding:"uuid"`
}

func (bc *BookBorrowValidator) Bind(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(bc)
}

func (bc *BookBorrowValidator) GetBorrowerID() string {
	return bc.BorrowerID
}
