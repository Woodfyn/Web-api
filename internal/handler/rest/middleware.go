package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"uri":    c.Request.RequestURI,
		}).Info()
	}
}
