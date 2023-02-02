package main

import (
	"fmt"
	"net/http"

	"jwt-go/initializers"

	"jwt-go/controllers"

	"github.com/gin-gonic/gin"
	"jwt-go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "we are in home"})
}

func main() {
	fmt.Println("Starting API")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/api", home)
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	router.GET("/logout",middleware.RequireAuth,controllers.Logout)
	router.Run(":8080")

}
