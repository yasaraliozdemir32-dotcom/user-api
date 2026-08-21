package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"user-api/config"
	"user-api/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func UploadFile(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Kullanici kimligi bulunamadi",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Dosya bulunamadi",
		})
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Uploads klasoru olusturulamadi",
		})
	}

	token := uuid.New().String()

	filePath := filepath.Join(
		"uploads",
		token+"_"+file.Filename,
	)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Dosya kaydedilemedi",
		})
	}

	newFile := models.File{
		FileName:   file.Filename,
		FilePath:   filePath,
		FileSize:   file.Size,
		ShareToken: token,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if err := config.DB.Create(&newFile).Error; err != nil {
		os.Remove(filePath)

		return c.Status(500).JSON(fiber.Map{
			"error": "Dosya bilgisi kaydedilemedi",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Dosya basariyla yuklendi",
		"file":    newFile,
		"link":    fmt.Sprintf("/download/%s", token),
	})
}

func DownloadFile(c fiber.Ctx) error {
	token := c.Params("token")

	var file models.File

	if err := config.DB.
		Where("share_token = ?", token).
		First(&file).Error; err != nil {

		return c.Status(404).JSON(fiber.Map{
			"error": "Dosya bulunamadi",
		})
	}

	if time.Now().After(file.ExpiresAt) {
		return c.Status(410).JSON(fiber.Map{
			"error": "Dosya linkinin suresi dolmus",
		})
	}

	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "Dosya sunucuda bulunamadi",
		})
	}

	return c.Download(file.FilePath, file.FileName)
}

func GetMyFiles(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Kullanici kimligi bulunamadi",
		})
	}

	var files []models.File

	if err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&files).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"error": "Dosyalar getirilemedi",
		})
	}

	return c.Status(200).JSON(files)
}

func DeleteFile(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Kullanici kimligi bulunamadi",
		})
	}

	id := c.Params("id")

	var file models.File

	if err := config.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&file).Error; err != nil {

		return c.Status(404).JSON(fiber.Map{
			"error": "Dosya bulunamadi",
		})
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return c.Status(500).JSON(fiber.Map{
			"error": "Dosya silinemedi",
		})
	}

	if err := config.DB.Delete(&file).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Dosya bilgisi silinemedi",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Dosya basariyla silindi",
	})
}