package controllers

import (
	"github.com/Aeroxee/kafekoding-api/handlers"
	"github.com/gin-gonic/gin"
)

func UserController(group *gin.RouterGroup) {
	userHandler := handlers.NewUserHandlerV1()
	group.GET("/auth", userHandler.CheckAuthHandler)
	group.PUT("/update-info", userHandler.UpdateInfoUserHandler)
	group.POST("/change-password", userHandler.ChangePasswordHandler)
}
