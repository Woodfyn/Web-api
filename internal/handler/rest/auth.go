package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/gin-gonic/gin"
)

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

	if err := h.services.SignUp(input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

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

	_, token, err := h.services.SignIn(input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(c, err)
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func handleNotFoundError(c *gin.Context, err error) {
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadRequest, string(response))
}
