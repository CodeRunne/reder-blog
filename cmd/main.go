package main

import (
	"os"
	
	_ "github.com/coderunne/jwt-login/pkgs/config"
	"github.com/coderunne/jwt-login/pkgs/routes"
)

func main() {

	var FRONTEND_URL = os.Getenv("FRONTEND_URL")
	
	// Define cors for the application
	origins := []string{FRONTEND_URL}
	routes.Cors(origins...)

	// Define application routes
	routes.GuestRoutes()
	routes.UserRoutes()
	routes.CategoryRoutes()
	routes.TagRoutes()
	routes.PostRoutes()
	routes.DashboardRoutes()
	routes.Start(":3000")
}
