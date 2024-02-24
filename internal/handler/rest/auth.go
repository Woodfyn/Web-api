package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.SignUpInput true "account info"
// @Success 200
// @Failure 400,500 {object} errorResponse
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

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,500 {object} errorResponse
// @Router /auth/sign-in [post]
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

	accessToken, refreshToken, err := h.services.Users.SignIn(input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(c, err)
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", "Authorization=Bearer "+refreshToken+"; HttpOnly; Path=/; Max-Age=3600")
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

// @Summary Refresh
// @Tags auth
// @Description refresh
// @ID refresh
// @Accept json
// @Produce json
// @Success 200 {string} string "token"
// @Failure 400,500 {object} errorResponse
// @Router /auth/refresh [get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken, err := getTokenFromCookie(cookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Info(cookie.Value)

	accessToken, refreshToken, err := h.services.Users.RefreshTokens(refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", "Authorization=Bearer "+refreshToken+"; HttpOnly; Path=/; Max-Age=3600")
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

func handleNotFoundError(c *gin.Context, err error) {
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, string(response))
}

func getTokenFromCookie(cookieValue string) (string, error) {
	headerParts := strings.Split(cookieValue, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

// @Summary LogOut
// @Tags auth
// @Description logout
// @ID logout
// @Accept json
// @Produce json
// @Success 200
// @Failure 400,500 {object} errorResponse
// @Router /auth/logout [post]
func (h *Handler) logOut(c *gin.Context) {
	cookie, err := c.Request.Cookie("Authorization")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken, err := getTokenFromCookie(cookie.Value)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Users.LogOut(refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	logrus.Info("User logged out")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, nil)
}
