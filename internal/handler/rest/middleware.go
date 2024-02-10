package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
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

func authMiddleware(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := getUserIdFromRequest(c, store)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
	}
}

func getUserIdFromRequest(c *gin.Context, store *sessions.CookieStore) error {
	session, err := store.Get(c.Request, "cookie-name")
	if err != nil {
		return err
	}

	_, ok := session.Values["user_id"].(int)
	if !ok {
		return errors.New("can't get user id from session")
	}

	return nil
}
