package rest

import (
	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type getAllGameResponse struct {
	Data []domain.Game `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
