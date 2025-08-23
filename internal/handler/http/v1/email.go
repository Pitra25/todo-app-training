package v1

import (
	"net/http"
	errorsResponse "todo-app/internal/errors"

	"todo-app/internal/repository/mysql/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) postSendEmail(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input models.SendEmailInput
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Emails.SendEmail(input.To, userId); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "email sent successfully",
	})
}

func (h *Handler) postConfirmationEmail(c *gin.Context) {
	userId, err := h.mw.GetUserId(c)
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input models.SendConfirmationEmailInput
	if err := c.BindJSON(&input); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Emails.ConfirmationEmail(input.Code, userId); err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "confirmation email sent successfully",
	})
}
