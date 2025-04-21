package setup

import (
	"github.com/csokviktor/lib_manager/internal/api/authors"
	"github.com/csokviktor/lib_manager/internal/api/books"
	"github.com/csokviktor/lib_manager/internal/api/borrowers"
	"github.com/csokviktor/lib_manager/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthorsRoutes(
	router *gin.RouterGroup,
	authorsService service.Authors,
) {
	routeHandler := authors.NewRoutes(authorsService)
	v1Group := router.Group("/v1/authors")

	v1Group.POST("", routeHandler.CreateAuthor)
	v1Group.GET("", routeHandler.GetAuthors)
	v1Group.GET("/:id", routeHandler.GetAuthor)
}

func BorrowersRoutes(
	router *gin.RouterGroup,
	borrowersService service.Borrowers,
) {
	routeHandler := borrowers.NewRoutes(borrowersService)
	v1Group := router.Group("/v1/borrowers")

	v1Group.POST("", routeHandler.CreateBorrower)
	v1Group.GET("", routeHandler.GetBorrowers)
	v1Group.GET("/:id", routeHandler.GetBorrower)
	v1Group.GET("/:id/books", routeHandler.GetBorrowerBooks)
}

func BooksRoutes(
	router *gin.RouterGroup,
	booksService service.Books,
) {
	routeHandler := books.NewRoutes(booksService)
	v1Group := router.Group("/v1/books")

	v1Group.POST("", routeHandler.CreateBook)
	v1Group.GET("", routeHandler.GetBooks)
	v1Group.POST("/:id/borrow", routeHandler.BorrowBook)
}
