package controllers

import (
	"github.com/Aeroxee/kafekoding-api/handlers"
	"github.com/gin-gonic/gin"
)

func ClassControllerV1WithAuth(group *gin.RouterGroup) {
	classHandlerV1 := handlers.NewClassHandlerV1()

	group.POST("", classHandlerV1.CreateHandler)
	group.GET("/:slug", classHandlerV1.Detail)
	group.PUT("/:slug", classHandlerV1.Update)
	group.DELETE("/:slug", classHandlerV1.Delete)
}

func ClassControllerV1NoAuth(group *gin.RouterGroup) {
	classHandlerV1 := handlers.NewClassHandlerV1()
	group.GET("", classHandlerV1.Get)
}
