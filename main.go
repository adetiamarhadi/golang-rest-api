package main

import (
	"github.com/adetiamarhadi/golang-rest-api/config"
	"github.com/adetiamarhadi/golang-rest-api/controller"
	"github.com/adetiamarhadi/golang-rest-api/middleware"
	"github.com/adetiamarhadi/golang-rest-api/repository"
	"github.com/adetiamarhadi/golang-rest-api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()
var userRepository repository.UserRepository = repository.NewUserRepository(db)
var bookRepository repository.BookRepository = repository.NewBookRepository(db)
var jwtService service.JWTService = service.NewJWTService()
var authService service.AuthService = service.NewAuthService(userRepository)
var userService service.UserService = service.NewUserService(userRepository)
var bookService service.BookService = service.NewBookService(bookRepository)
var authController controller.AuthController = controller.NewAuthController(authService, jwtService)
var userController controller.UserController = controller.NewUserController(userService, jwtService)
var bookController controller.BookController = controller.NewBookController(bookService, jwtService)

func main() {
	defer config.CloseDatabaseConnection(db)

	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	r.Run()
}
