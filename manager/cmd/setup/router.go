package setup

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RouterEngine() *gin.Engine {
	// setup basic router settings
	router := gin.New()
	router.Use(ZerologMiddleware(log.Logger))
	router.Use(TimeoutMiddleware(10 * time.Second))
	router.Use(gin.Recovery())
	// used only for the demo project, must be configured if deployed
	router.Use(cors.New(cors.Config{AllowAllOrigins: true}))

	return router
}
