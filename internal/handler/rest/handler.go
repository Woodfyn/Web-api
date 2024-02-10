package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Woodfyn/Web-api/docs"
	"github.com/Woodfyn/Web-api/internal/service"
)

type Handler struct {
	services       *service.Services
	cookieSessions *sessions.CookieStore
}

func NewHandler(services *service.Services, cookieSessions *sessions.CookieStore) *Handler {
	return &Handler{
		services:       services,
		cookieSessions: cookieSessions,
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
		auth.GET("/log-out", h.logOut)
	}

	api := router.Group("/api")
	{
		game := api.Group("/game")
		{
			game.Use(authMiddleware(h.cookieSessions))

			game.POST("/", h.addGame)
			game.GET("/", h.getAllGame)
			game.GET("/:id", h.getGameByID)
			game.PUT("/:id", h.updateGameByID)
			game.DELETE("/:id", h.deleteGameByID)
		}

	}

	return router
}
