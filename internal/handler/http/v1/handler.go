package v1

import (
	"todo-app/internal/handler/http"
	"todo-app/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "todo-app/docs"
)

type Handler struct {
	services *service.Service
	mw       *http.MW
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.singUp)
		auth.POST("/sing-in", h.singIp)
	}

	api := router.Group("/api", http.NewMW(h.services).UserIdentity)
	{
		user := api.Group("users")
		{
			user.GET("/", h.getAllUsers)
			user.GET("/:id", h.getUserById)
			user.PUT("/:id", h.updateUser)
			user.DELETE("/delete", h.deleteAccount)
		}

		list := api.Group("lists")
		{
			list.POST("/", h.createList)
			list.GET("/", h.getAllLists)
			list.GET("/:id", h.getListsById)
			list.PUT("/:id", h.updateList)
			list.DELETE("/:id", h.deleteList)

			items := list.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItemsList)
			}
		}

		items := api.Group("items")
		{
			items.GET("/", h.getAllItem)
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}

		email := api.Group("email")
		{
			email.POST("/send", h.postSendEmail)
			email.POST("/confirmation", h.postConfirmationEmail)
		}
	}

	return router
}
