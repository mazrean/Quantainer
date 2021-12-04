package gorm2

import (
	"time"

	"github.com/google/uuid"
)

var (
	tables = []interface{}{}
)

type FileTable struct {
	ID         uuid.UUID     `gorm:"type:varchar(36);not null;primaryKey"`
	FileTypeID int           `gorm:"type:tinyint;not null"`
	CreatedAt  time.Time     `gorm:"type:datetime;not null"`
	FileType   FileTypeTable `gorm:"foreignKey:FileTypeID"`
}

func (ft *FileTable) TableName() string {
	return "game_files"
}

type FileTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (ft *FileTypeTable) TableName() string {
	return "file_types"
}
