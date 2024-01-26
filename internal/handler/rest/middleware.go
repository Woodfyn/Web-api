package rest

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL)
	}
}
