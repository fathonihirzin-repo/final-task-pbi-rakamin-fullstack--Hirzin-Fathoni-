package main

import (
	"crud/controllers"
	"crud/initializers"
	"crud/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequiredAuth, controllers.Validate)
	r.GET("/users", controllers.UserIndex)
	r.GET("/user/:id", controllers.GetUser)
	r.PUT("user/:id", controllers.UpdateUser)
	r.DELETE("user/:id", controllers.DeleteUser)

	r.Run()
}
