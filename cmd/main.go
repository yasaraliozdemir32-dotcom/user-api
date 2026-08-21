package main

import (
	"log"
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

	log.Fatal(app.Listen(":8082"))
}