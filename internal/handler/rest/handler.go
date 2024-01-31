package rest

import (
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Woodfyn/Web-api/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(loggingMiddleware())

	api := router.Group("/api")
	{
		game := api.Group("/game")
		{
			game.POST("/", h.addGame)
			game.GET("/", h.getAllGame)
			game.GET("/:id", h.getGameByID)
			game.PUT("/:id", h.updateGameByID)
			game.DELETE("/:id", h.deleteGameByID)
		}

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	}

	return router
}
