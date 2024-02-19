package controllers

import (
	"github.com/Aeroxee/kafekoding-api/handlers"
	"github.com/gin-gonic/gin"
)

func ArticleControllerNoAuth(group *gin.RouterGroup) {
	articleHandlerV1 := handlers.NewArticleHandlerV1()

	group.GET("", articleHandlerV1.Get)
}

func ArticleControllerWithAuth(group *gin.RouterGroup) {
	articleHandlerV1 := handlers.NewArticleHandlerV1()
	group.POST("", articleHandlerV1.CreateHandler)
}
