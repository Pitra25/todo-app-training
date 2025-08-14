package v1

import (
	"net/http"
	"todo-app/internal/errors"
	"todo-app/internal/repository/mysql/models"

	"github.com/gin-gonic/gin"
)

// @Summary Get user By Id
// @Security ApiKeyAuth
// @Tags Users
// @Description get user by id
// @ID get-user-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} modules.Users
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/:id [get]
func (h *Heandler) getUserById(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.Users.GetUserById(userId)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete User Account
// @Security ApiKeyAuth
// @Tags Ssers
// @Description delete current user's account
// @ID delete-user-account
// @Accept  json
// @Produce  json
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/delete [delete]
func (h *Heandler) deleteAccount(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Users.DeleteUser(userId)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{
		Status: "ok",
	})
}

// @Summary Update User Account
// @Security ApiKeyAuth
// @Tags Users
// @Description update current user's account
// @ID update-user-account
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param input body models.UpdadeUsers true "User update input"
// @Success 200 {object} errorsResponse.statusResponse
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/users/{id} [put]
func (h *Heandler) updateUser(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input *models.UpdateUserInpur
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Users.UpdateUser(userId, input)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, errorsResponse.StatusResponse{
		Status: "ok",
	})
}

// @Summary Get All User list
// @Security ApiKeyAuth
// @Tags Users
// @Description get all user
// @ID get-all-item
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Users
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items [get]
func (h *Heandler) getAllUsers(c *gin.Context) {
	users, err := h.services.Users.GetUserAll()
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}
