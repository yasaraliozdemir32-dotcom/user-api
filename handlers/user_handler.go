package handlers

import (
	"user-api/models"
	"user-api/service"

	"github.com/gofiber/fiber/v3"
)

func CreateUser(c fiber.Ctx) error {
	var user models.User

	if err := c.Bind().Body(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Gecersiz veri",
		})
	}

	if err := service.CreateUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Kullanici olusturulamadi",
		})
	}

	return c.Status(201).JSON(user)
}

func GetAllUsers(c fiber.Ctx) error {
	users, err := service.GetAllUsers()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Kullanicilar getirilemedi",
		})
	}

	return c.Status(200).JSON(users)
}

func GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")

	user, err := service.GetUserByID(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Kullanici bulunamadi",
		})
	}

	return c.Status(200).JSON(user)
}

func UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")

	user, err := service.GetUserByID(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Kullanici bulunamadi",
		})
	}

	var updateData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Gecersiz veri",
		})
	}

	if updateData.Name != "" {
		user.Name = updateData.Name
	}

	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	if updateData.Password != "" {
		user.Password = updateData.Password
	}

	if err := service.UpdateUser(&user, updateData.Password != ""); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Kullanici guncellenemedi",
		})
	}

	return c.Status(200).JSON(user)
}

func DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	user, err := service.GetUserByID(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Kullanici bulunamadi",
		})
	}

	if err := service.DeleteUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Kullanici silinemedi",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Kullanici basariyla silindi",
	})
}

func Login(c fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&loginData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Gecersiz veri",
		})
	}

	token, err := service.Login(
		loginData.Email,
		loginData.Password,
	)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Email veya sifre hatali",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Giris basarili",
		"token":   token,
	})
}
