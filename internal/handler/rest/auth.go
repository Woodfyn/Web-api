package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/gin-gonic/gin"
)

const cookieName = "cookie-name"

// @Summary signUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.SignUpInput true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,500 {object} errorResponce
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.SignUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Users.SignUp(input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary signIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {string} string "session"
// @Failure 400,500 {object} errorResponce
// @Router /auth/sign-in [get]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.cookieSessions.Get(c.Request, cookieName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	sessionWithValue, err := h.services.Users.SignIn(input, session)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(c, err)
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.cookieSessions.Save(c.Request, c.Writer, sessionWithValue)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary logout
// @Tags auth
// @Description logout
// @ID log-out
// @Accept json
// @Produce json
// @Success 200 {string} string "session"
// @Failure 400,500 {object} errorResponce
// @Router /auth/log-out [get]
func (h *Handler) logOut(c *gin.Context) {
	session, err := h.cookieSessions.Get(c.Request, cookieName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sessionWithCacnel, err := h.services.Users.LogOut(session)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.cookieSessions.Save(c.Request, c.Writer, sessionWithCacnel)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func handleNotFoundError(c *gin.Context, err error) {
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, string(response))
}
