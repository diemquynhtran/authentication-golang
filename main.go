package main

import (
	"jwt-gin/config"
	"jwt-gin/controller"
	"jwt-gin/middleware"
	"jwt-gin/repository"
	"jwt-gin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUSerRepository(db)
	jwtService service.JWTService = service.NewJWTService()
	authService service.AuthService = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userSerivice  service.UserService = service.NewUSerRepository(userRepository)
	userController controller.UserController = controller.NewUserController(userSerivice, jwtService)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	bookService service.BookService = service.NewBookService(bookRepository)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)


)


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
		userRoutes.PUT("/update", userController.Update)
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