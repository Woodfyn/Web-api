package rest

import (
	"net/http"
	"strconv"

	"github.com/Woodfyn/Web-api/internal/domain"

	"github.com/gin-gonic/gin"
)

// @Summary AddGame
// @Tags game
// @Description add new game
// @ID add-game
// @Accept json
// @Produce json
// @Param input body domain.Game true "game info"
// @Success 200 {integer} integer 1
// @Failure 400,500 {object} errorResponce
// @Router /game [post]
func (h *Handler) addGame(c *gin.Context) {
	var input domain.Game
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Game.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary GetAll
// @Tags game
// @Description getAll games
// @ID get-all
// @Accept json
// @Produce json
// @Success 200 {object} getAllGameResponse
// @Failure 400,500 {object} errorResponce
// @Router /game [get]
func (h *Handler) getAllGame(c *gin.Context) {
	games, err := h.services.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllGameResponse{
		Data: games,
	})
}

// @Summary GetGameByID
// @Tags game
// @Description get game by id
// @ID get-game-by-id
// @Accept json
// @Produce json
// @Param id path int true "game id"
// @Success 200 {object} domain.Game
// @Failure 400,500 {object} errorResponce
// @Router /game/{id} [get]
func (h *Handler) getGameByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	game, err := h.services.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, game)
}

// @Summary UpdateGameByID
// @Tags game
// @Description update game by id
// @ID update-game
// @Accept json
// @Produce json
// @Param id path int true "game id"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponce
// @Router /game/{id} [put]
func (h *Handler) updateGameByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input domain.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UpdateById(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary DeleteGameByID
// @Tags game
// @Description delete game by id
// @ID delete-game
// @Accept json
// @Produce json
// @Param id path int true "game id"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponce
// @Router /game/{id} [delete]
func (h *Handler) deleteGameByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.DeleteById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
