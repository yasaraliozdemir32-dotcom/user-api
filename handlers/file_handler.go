package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"user-api/config"
	"user-api/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// =====================================================
// DOSYA YUKLE
// =====================================================

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

	// =================================================
	// 30 MB DOSYA BOYUTU SINIRI
	// =================================================

	const maxFileSize = 30 * 1024 * 1024

	if file.Size > maxFileSize {
		return c.Status(413).JSON(fiber.Map{
			"error": "Dosya boyutu 30 MB'dan buyuk olamaz",
		})
	}

	// =================================================
	// DOSYA TURU KONTROLU
	// =================================================

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".zip":  true,
		".rar":  true,
		".txt":  true,
	}

	extension := strings.ToLower(
		filepath.Ext(file.Filename),
	)

	if !allowedExtensions[extension] {
		return c.Status(415).JSON(fiber.Map{
			"error": "Bu dosya turune izin verilmiyor",
		})
	}

	// =================================================
	// UPLOADS KLASORU
	// =================================================

	if err := os.MkdirAll("uploads", 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Uploads klasoru olusturulamadi",
		})
	}

	// =================================================
	// DOSYA TOKEN
	// =================================================

	token := uuid.New().String()

	filePath := filepath.Join(
		"uploads",
		token+"_"+file.Filename,
	)

	// =================================================
	// DOSYAYI KAYDET
	// =================================================

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Dosya kaydedilemedi",
		})
	}

	// =================================================
	// VERITABANINA KAYDET
	// =================================================

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

	// =================================================
	// BASARILI
	// =================================================

	return c.Status(201).JSON(fiber.Map{
		"message": "Dosya basariyla yuklendi",
		"file":    newFile,
		"link":    fmt.Sprintf("/download/%s", token),
	})
}

// =====================================================
// DOSYA INDIR
// =====================================================

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

	return c.Download(
		file.FilePath,
		file.FileName,
	)
}

// =====================================================
// KULLANICININ DOSYALARINI GETIR
// =====================================================

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

// =====================================================
// DOSYA SIL
// =====================================================

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

	if err := os.Remove(file.FilePath); err != nil &&
		!os.IsNotExist(err) {

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
