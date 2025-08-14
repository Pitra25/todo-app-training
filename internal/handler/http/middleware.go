package http

import (
	"errors"
	"net/http"
	"strings"
	errorsResponse "todo-app/internal/errors"
	"todo-app/internal/service"

	"github.com/gin-gonic/gin"
)

type MW struct {
	services *service.Service
}

func NewMW(services *service.Service) *MW {
	return &MW{services: services}
}

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (mw *MW) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		errorsResponse.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		errorsResponse.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := mw.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		errorsResponse.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func (mw *MW) GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		errorsResponse.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
