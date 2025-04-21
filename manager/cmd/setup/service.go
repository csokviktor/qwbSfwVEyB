package setup

import (
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/service"
	"gorm.io/gorm"
)

func AuthorsService(db *gorm.DB) service.Authors {
	rep := repository.NewAuthors(db)
	return service.NewAuthors(rep)
}

func BorrowersService(db *gorm.DB) service.Borrowers {
	rep := repository.NewBorrowers(db)
	return service.NewBorrowers(rep)
}

func BooksService(db *gorm.DB, authorsService service.Authors, borrowersService service.Borrowers) service.Books {
	rep := repository.NewBooks(db)
	return service.NewBooks(rep, authorsService, borrowersService)
}
