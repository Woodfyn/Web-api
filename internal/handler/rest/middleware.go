package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"uri":    c.Request.RequestURI,
		}).Info()

		c.Next()
	}
}

func (h *Handler) userIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c.Request)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := h.services.Users.ParseToken(token)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		logrus.Infof("User with id %s auth", userId)

		c.Next()
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
