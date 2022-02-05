package main

import (
	"github.com/adetiamarhadi/golang-rest-api/config"
	"github.com/adetiamarhadi/golang-rest-api/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()
var authController controller.AuthController = controller.NewAuthController()

func main() {
	defer config.CloseDatabaseConnection(db)

	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	r.Run()
}
