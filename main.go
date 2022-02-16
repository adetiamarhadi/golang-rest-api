package main

import (
	"github.com/adetiamarhadi/golang-rest-api/config"
	"github.com/adetiamarhadi/golang-rest-api/controller"
	"github.com/adetiamarhadi/golang-rest-api/repository"
	"github.com/adetiamarhadi/golang-rest-api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()
var userRepository repository.UserRepository = repository.NewUserRepository(db)
var jwtService service.JWTService = service.NewJWTService()
var authService service.AuthService = service.NewAuthService(userRepository)
var authController controller.AuthController = controller.NewAuthController(authService, jwtService)

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
