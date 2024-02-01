package routes

import (
	"fmt"
	"log"

	"strings"

	"github.com/coderunne/jwt-login/pkgs/controllers"
	"github.com/coderunne/jwt-login/pkgs/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func init() {
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 4 << 20 // 4 MiB
}

func Cors(origins ...string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowMethods = []string{"PUT", "DELETE", "POST", "GET"}
	config.AllowCredentials = true

	router.Use(cors.New(config))
}

func GuestRoutes() {
	router.POST("/register", controllers.RegisterController)
	router.POST("/login", controllers.LoginController)
}

func UserRoutes() {
	users := router.Group("/users")
	// Middleware
	users.Use(middleware.AuthenticateSession())
	// Unauthenticated User Route
	router.GET("/users", controllers.UsersController)
	// Authenticated User Routes
	users.GET("/:id", controllers.UserGetController)
	users.GET("/:id/posts", controllers.UserPostsController)
}

func CategoryRoutes() {
	categories := router.Group("/categories")
	// Middleware
	categories.Use(middleware.AuthenticateSession())
	// Unauthenticated Category Routes
	router.GET("/categories", controllers.CategoriesAllController)
	router.GET("/categories/:id/posts", controllers.CategoryPostsController)
	// Authenticated Posts Routes
	categories.POST("/create", controllers.CategoriesCreateController)
	categories.PUT("/:id/update", controllers.CategoriesUpdateController)
	categories.DELETE("/:id/delete", controllers.CategoriesDeleteController)
}

func TagRoutes() {
	tags := router.Group("/tags")
	// Middleware
	tags.Use(middleware.AuthenticateSession())
	// Unauthenticated Category Routes
	router.GET("/tags", controllers.TagsAllController)
	router.GET("/tags/:id/posts", controllers.TagPostsController)
	// Authenticated Posts Routes
	tags.POST("/create", controllers.TagCreateController)
	tags.PUT("/:id/update", controllers.TagUpdateController)
	tags.DELETE("/:id/delete", controllers.TagDeleteController)
}

func PostRoutes() {
	posts := router.Group("/posts")
	// Middleware
	posts.Use(middleware.AuthenticateSession())
	// Unauthenticated Post Route
	router.GET("/posts", controllers.PostAllController)
	router.GET("/posts/:slug", controllers.PostGetController)
	// Authenticated Posts Routes
	posts.POST("/create", controllers.PostCreateController)
	posts.PUT("/:slug/update", controllers.PostUpdateController)
	posts.DELETE("/:slug/delete", controllers.PostDeleteController)
}

func DashboardRoutes() {
	dashboard := router.Group("/dashboard")
	// Middleware
	dashboard.Use(middleware.AuthenticateSession())
	// Dashboard Routes
	dashboard.GET("/", controllers.DashboardController)
}

func Start(port string) {
	if port == "" {
		log.Fatal("Provide a port to start server!")
	} else if strings.Contains(port, ":") == false {
		port = fmt.Sprintf(":%s", port)
	}
	router.Run(port)
}
