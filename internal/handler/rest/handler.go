package rest

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Woodfyn/Web-api/docs"
	"github.com/Woodfyn/Web-api/internal/service"
	"github.com/Woodfyn/Web-api/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(loggingMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

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
