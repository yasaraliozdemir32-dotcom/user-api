package main

import (
	"log"
	"os"
	"user-api/config"
	"user-api/models"
	"user-api/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(
		&models.User{},
		&models.File{},
	)

	app := fiber.New()

	routes.UserRoutes(app)
	app.Use("/", static.New("./web"))

	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
	}

	log.Fatal(app.Listen(":" + port))
}
