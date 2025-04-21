package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/csokviktor/qwbSfwVEyB/manager/cmd/setup"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := setup.Config()
	if err != nil {
		log.Fatal().Msgf("startup failed while reading config: %v", err)
	}

	setup.Logger(cfg)

	// Setup Database
	if err = setup.RunMigrationScripts(cfg); err != nil {
		log.Fatal().Err(err).Msg("DB migration failed")
	}
	db, err := setup.ConnectToDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to DB failed")
	}

	router := setup.RouterEngine()

	// Setup service level
	authorsService := setup.AuthorsService(db)
	borrowersService := setup.BorrowersService(db)
	booksService := setup.BooksService(db, authorsService, borrowersService)

	// Setup route handlers
	setup.AuthorsRoutes(&router.RouterGroup, authorsService)
	setup.BorrowersRoutes(&router.RouterGroup, borrowersService)
	setup.BooksRoutes(&router.RouterGroup, booksService)

	// Run server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msgf("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Info().Msgf("Server Shutdown: %s", err)
	}

	log.Info().Msgf("Server exiting")
}
