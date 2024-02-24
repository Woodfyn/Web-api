package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Woodfyn/Web-api/docs"
	"github.com/Woodfyn/Web-api/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(loggingMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
		auth.POST("/logout", h.logOut)
	}

	api := router.Group("/api", h.userIdentity())
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

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}
