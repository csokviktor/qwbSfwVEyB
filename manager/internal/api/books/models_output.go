package books

import "github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"

type APIBook struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	AuthorID   string  `json:"authorID"`
	BorrowerID *string `json:"borrowerID,omitempty"`
}

func DBToAPI(dbBook *dbmodels.Book) *APIBook {
	return &APIBook{
		ID:         dbBook.ID,
		Title:      dbBook.Title,
		AuthorID:   dbBook.AuthorID,
		BorrowerID: dbBook.BorrowerID,
	}
}

func DBsToAPIs(dbBooks []dbmodels.Book) []APIBook {
	apiBook := make([]APIBook, 0, len(dbBooks))

	for _, dbBook := range dbBooks {
		apiBook = append(apiBook, *DBToAPI(&dbBook))
	}
	return apiBook
}
