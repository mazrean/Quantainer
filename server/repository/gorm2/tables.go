package gorm2

import (
	"time"

	"github.com/google/uuid"
)

var (
	tables = []interface{}{
		&FileTable{},
		&FileTypeTable{},
		&ResourceTable{},
	}
)

type FileTable struct {
	ID         uuid.UUID     `gorm:"type:varchar(36);not null;primaryKey"`
	FileTypeID int           `gorm:"type:tinyint;not null"`
	CreatorID  uuid.UUID     `gorm:"type:varchar(36);not null"`
	CreatedAt  time.Time     `gorm:"type:datetime;not null"`
	FileType   FileTypeTable `gorm:"foreignKey:FileTypeID"`
}

func (ft *FileTable) TableName() string {
	return "files"
}

type FileTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (ft *FileTypeTable) TableName() string {
	return "file_types"
}

type ResourceTable struct {
	ID             uuid.UUID         `gorm:"type:varchar(36);not null;primaryKey"`
	FileID         uuid.UUID         `gorm:"type:varchar(36);not null"`
	Name           string            `gorm:"type:varchar(64);size:64;not null"`
	ResourceTypeID int               `gorm:"type:tinyint;not null"`
	Comment        string            `gorm:"type:varchar(400);size:400;not null"`
	CreatedAt      time.Time         `gorm:"type:datetime;not null"`
	File           FileTable         `gorm:"foreignKey:FileID"`
	ResourceType   ResourceTypeTable `gorm:"foreignKey:ResourceTypeID"`
}

func (rt *ResourceTable) TableName() string {
	return "resources"
}

type ResourceTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (rtt *ResourceTypeTable) TableName() string {
	return "resource_types"
}
