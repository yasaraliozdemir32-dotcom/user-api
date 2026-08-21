package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FileName    string         `json:"file_name"`
	FilePath    string         `json:"-"`
	FileSize    int64          `json:"file_size"`
	ShareToken  string         `json:"share_token"`
	UserID      uint           `json:"user_id"`
	ExpiresAt   time.Time      `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}