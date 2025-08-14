package v1

import (
	"net/http"
	"strconv"
	errorsResponse "todo-app/internal/errors"
	"todo-app/internal/repository/mysql/models"

	"github.com/gin-gonic/gin"
)

// @Summary Create todo item
// @Security ApiKeyAuth
// @Tags Items
// @Description create todo item
// @ID create-item
// @Accept  json
// @Produce  json
// @Param input body types.TodoItems true "item info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/item [post]
func (h *Heandler) createItem(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input models.TodoItems
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItems.Create(userId, listId, input)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get All items list
// @Security ApiKeyAuth
// @Tags Items
// @Description get all items
// @ID get-all-item
// @Accept  json
// @Produce  json
// @Success 200 {object} models.TodoItems
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items [get]
func (h *Heandler) getAllItemsList(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	item, err := h.services.TodoItems.GetAllItemsList(userId, listId)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)

}

// @Summary Get All items list
// @Security ApiKeyAuth
// @Tags Items
// @Description get all items
// @ID get-all-item
// @Accept  json
// @Produce  json
// @Success 200 {object} models.TodoItems
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items [get]
func (h *Heandler) getAllItem(c *gin.Context) {

	item, err := h.services.TodoItems.GetAllItem()
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)

}

// @Summary Get item By Id
// @Security ApiKeyAuth
// @Tags Items
// @Description get items by id
// @ID get-items-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.TodoItems
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [get]
func (h *Heandler) getItemById(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	item, err := h.services.TodoItems.GetById(userId, id)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Update Item
// @Security ApiKeyAuth
// @Tags Items
// @Description update Item by id
// @ID update-Item
// @Accept  json
// @Produce  json
// @Param id path int true "Item ID"
// @Param input body models.UpdadeListInput true "Item update input"
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 400 {object} errorResponse "Invalid ID or input"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/items/{id} [put]
func (h *Heandler) updateItem(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdadeItemInput
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoItems.Update(userId, id, input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{Status: "ok"})
}

// @Summary Delete Item
// @Security ApiKeyAuth
// @Tags Items
// @Description delete Item by id
// @ID delete-Item
// @Accept  json
// @Produce  json
// @Param id path int true "Item ID"
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 400 {object} errorResponse "Invalid ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/items/{id} [delete]
func (h *Heandler) deleteItem(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoItems.Delete(userId, id)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{
		Status: "ok",
	})
}
