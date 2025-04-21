package authors

import (
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/api/books"
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
)

type APIAuthor struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Books []books.APIBook `json:"books"`
}

func DBToAPI(dbAuthor *dbmodels.Author) *APIAuthor {
	return &APIAuthor{
		ID:    dbAuthor.ID,
		Name:  dbAuthor.Name,
		Books: books.DBsToAPIs(dbAuthor.Books),
	}
}

func DBsToAPIs(dbAuthors []dbmodels.Author) []APIAuthor {
	apiAuthors := make([]APIAuthor, 0, len(dbAuthors))

	for _, dbAuthor := range dbAuthors {
		apiAuthors = append(apiAuthors, *DBToAPI(&dbAuthor))
	}
	return apiAuthors
}
