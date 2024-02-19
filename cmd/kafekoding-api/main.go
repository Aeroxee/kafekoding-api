package main

import (
	"github.com/Aeroxee/kafekoding-api/controllers"
	"github.com/Aeroxee/kafekoding-api/handlers"
	"github.com/Aeroxee/kafekoding-api/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	r.Run(":8000")
}
