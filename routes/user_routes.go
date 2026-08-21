package routes

import (
	"user-api/handlers"
	"user-api/middleware"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {

	// Login açık
	app.Post("/login", handlers.Login)

	// Kullanıcı oluşturma açık
	app.Post("/users", handlers.CreateUser)

	// Token gerekli
	app.Get("/users", middleware.JWTProtected, handlers.GetAllUsers)

	app.Get("/users/:id", middleware.JWTProtected, handlers.GetUserByID)

	app.Put("/users/:id", middleware.JWTProtected, handlers.UpdateUser)

	app.Delete("/users/:id", middleware.JWTProtected, handlers.DeleteUser)

	app.Get("/files", middleware.JWTProtected, handlers.GetMyFiles)

	app.Delete("/files/:id", middleware.JWTProtected, handlers.DeleteFile)

	// Dosya upload
	app.Post("/upload", middleware.JWTProtected, handlers.UploadFile)

	// Dosya download
	app.Get("/download/:token", handlers.DownloadFile)
}