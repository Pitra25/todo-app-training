package v1

import (
	"net/http"
	"strconv"
	errorsResponse "todo-app/internal/errors"
	"todo-app/internal/repository/mysql/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body models.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	var input models.TodoList
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.TodoList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lists)
}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.TodoList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [get]
func (h *Handler) getListsById(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Update List
// @Security ApiKeyAuth
// @Tags lists
// @Description update list by id
// @ID update-list
// @Accept  json
// @Produce  json
// @Param id path int true "List ID"
// @Param input body models.UpdateListInput true "List update input"
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 400 {object} errorResponse "Invalid ID or input"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoList.Update(userId, id, input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{Status: "ok"})
}

// @Summary Delete List
// @Security ApiKeyAuth
// @Tags lists
// @Description delete list by id
// @ID delete-list
// @Accept  json
// @Produce  json
// @Param id path int true "List ID"
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 400 {object} errorResponse "Invalid ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{
		Status: "ok",
	})
}
