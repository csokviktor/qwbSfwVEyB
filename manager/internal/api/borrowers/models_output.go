package borrowers

import (
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/api/books"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
)

type APIBorrower struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Books []books.APIBook `json:"books"`
}

func DBToAPI(dbAuthor *dbmodels.Borrower) *APIBorrower {
	return &APIBorrower{
		ID:    dbAuthor.ID,
		Name:  dbAuthor.Name,
		Books: books.DBsToAPIs(dbAuthor.Books),
	}
}

func DBsToAPIs(dbBorrowers []dbmodels.Borrower) []APIBorrower {
	apiBorrowers := make([]APIBorrower, 0, len(dbBorrowers))

	for _, dbBorrower := range dbBorrowers {
		apiBorrowers = append(apiBorrowers, *DBToAPI(&dbBorrower))
	}
	return apiBorrowers
}
