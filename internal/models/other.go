package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Permit int8

const (
	Read   = 1
	Write  = 2
	Manage = 3
)

type Share struct {
	ID           string    `gorm:"primaryKey;type:char(36)"`
	UserID       string    `gorm:"type:char(36)"`
	ResourceID   string    `gorm:"type:char(36)"`
	User         User      `gorm:"foreignKey:UserID"`
	Resource     Resource  `gorm:"foreignKey:ResourceID"`
	Token        string    `gorm:"size:64;unique;not null"` //分享的短链接标识
	PasswordHash string    `gorm:"size:255"`
	ExpiredAt    time.Time `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime:false"`
}

type Permission struct {
	ID         string    `gorm:"primaryKey;type:char(36)"`
	UserID     string    `gorm:"type:char(36)"`
	ResourceID string    `gorm:"type:char(36)"`
	User       User      `gorm:"foreignKey:UserID"`
	Resource   Resource  `gorm:"foreignKey:ResourceID"`
	Permit     int8      `gorm:"type:int(11)"`
	CreateAt   time.Time `gorm:"autoCreateTime:false"`
}

type Trash struct {
	ID               string    `gorm:"primaryKey;type:char(36)"`
	UserID           string    `gorm:"type:char(36)"`
	Resource         Resource  `gorm:"foreignKey:ResourceID"`
	ResourceID       string    `gorm:"type:char(36)"`
	User             User      `gorm:"foreignKey:UserID"`
	OriginalParentId string    `gorm:"type:char(36)"`
	DeleteTime       time.Time `gorm:"not null"`
	ExpireAt         time.Time `gorm:"not null"` //默认回收站自动清理的时间（天）
}

func (r *Trash) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	r.DeleteTime = time.Now()
	r.ExpireAt = time.Now().AddDate(0, 0, 30)
	return
}
