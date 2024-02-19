package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Aeroxee/kafekoding-api/auth"
	"github.com/Aeroxee/kafekoding-api/controllers"
	"github.com/Aeroxee/kafekoding-api/handlers"
	"github.com/Aeroxee/kafekoding-api/middlewares"
	"github.com/Aeroxee/kafekoding-api/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Static("/media", "./media")

	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}
	c := cors.New(config)
	r.Use(c)

	v1 := r.Group("/v1")

	// register
	userHandler := handlers.NewUserHandlerV1()
	v1.POST("/register", userHandler.RegisterHandler)
	v1.GET("/activate/:activationCode", userHandler.ActivationHandler)
	v1.POST("/get-token", userHandler.GetTokenHandler)

	userGroup := v1.Group("/user")
	userGroup.Use(middlewares.Authentication())
	controllers.UserController(userGroup)

	classGroupV1WithAuth := v1.Group("/classes")
	classGroupV1WithAuth.Use(middlewares.Authentication())
	controllers.ClassControllerV1WithAuth(classGroupV1WithAuth)

	classGroupV1NoAuth := v1.Group("/classes")
	controllers.ClassControllerV1NoAuth(classGroupV1NoAuth)

	// article group no auth
	articleGroupNoAuth := v1.Group("/articles")
	controllers.ArticleControllerNoAuth(articleGroupNoAuth)

	// article group with auth
	articleGroupWithAuth := v1.Group("/articles")
	articleGroupWithAuth.Use(middlewares.Authentication())
	controllers.ArticleControllerWithAuth(articleGroupWithAuth)

	// upload handler
	r.POST("/upload", func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		claims, err := auth.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		_, err = models.GetUserByID(claims.Credential.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Authentication is required.",
			})
			return
		}

		_, h, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		filename := h.Filename
		filenameUUID := uuid.NewString() + filepath.Ext(filename)
		destination := fmt.Sprintf("media/upload/%s", filenameUUID)

		err = ctx.SaveUploadedFile(h, destination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Upload a file successfully",
			"file":    destination,
			"ext":     filepath.Ext(filename),
		})
	})

	r.Run(":8000")
}
