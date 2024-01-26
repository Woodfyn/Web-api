package rest

import (
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/gin-gonic/gin"
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
	}

	return router
}
