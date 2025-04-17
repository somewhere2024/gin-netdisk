package models

import (
	"time"
)

type ResourceType int8

const (
	File   ResourceType = 1
	folder ResourceType = 2
)

type Resource struct {
	ID        string       `gorm:"primaryKey;type:char(36)"`
	UserId    string       `gorm:"type:char(36)"`
	Name      string       `gorm:"size:255;not null"`
	Type      ResourceType `gorm:"type:int(11)"`      //定义文件和文件夹类型，这里使用了自定义枚举
	ParentId  string       `gorm:"type:char(36)"`     //父文件夹id
	Size      int64        `gorm:"default:0"`         //文件大小；文件夹默认为零NULL
	MimeType  string       `gorm:"size:255;not null"` //文件类型
	Path      string       `gorm:"size:255;not null"`
	IsDelete  bool         `gorm:"default:false"`
	CreatedAt time.Time    `gorm:"not null"`
	UpdatedAt time.Time
}
