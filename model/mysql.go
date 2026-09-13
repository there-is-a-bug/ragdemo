package model

import (
	"time"

	"gorm.io/gorm"
)

type UserDocument struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64          `gorm:"column:user_id"`
	DocID     string         `gorm:"column:doc_id"`
	Title     string         `gorm:"column:title"`
	Content   string         `gorm:"column:content"`
	DocType   string         `gorm:"column:doc_type"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (u *UserDocument) TableName() string {
	return "document"
}
