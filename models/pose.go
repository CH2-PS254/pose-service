package models

import (
	"time"

	"gorm.io/gorm"
)

type Pose struct {
	ID          uint           `json:"id" gorm:"primaryKey" `
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CreatePoseInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdatePoseInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
