package main

import (
	"log"
	"os"
	"user-api/config"
	"user-api/models"
	"user-api/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(
	&models.User{},
	&models.File{},
)

	app := fiber.New() 

	routes.UserRoutes(app)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("API çalişiyor")
	})

	port := os.Getenv("PORT")
if port == "" {
	port = "8082"
}

log.Fatal(app.Listen(":" + port))
}