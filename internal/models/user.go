package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            string    `gorm:"primaryKey;type:char(36)"`
	Username      string    `gorm:"unique;size:50;not null"`
	Password_hash string    `gorm:"size:255;not null"` //保存加密后的密码
	Email         string    `gorm:"type:char(100);unique;not null"`
	StorageUsed   int64     `gorm:"default:0"`     //已存储的单位
	StorageTotal  int64     `gorm:"default:52428"` //总的存储空间，默认5GB
	CreatedAt     time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:true"`
}

type UserResponse struct {
	Username     string
	Email        string
	StorageUsed  int64
	StorageTotal int64
	CreatedAt    time.Time
	UpdateAt     time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return nil
}

type UserDetail struct {
	gorm.Model
	ID    string `gorm:"primaryKey;type:char(36)"`
	Hobby string `gorm:"type:char(100);unique;not null"`
}
