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
		&ResourceTypeTable{},
		&GroupTable{},
		&GroupTypeTable{},
		&ReadPermissionTable{},
		&WritePermissionTable{},
		&AdministratorTable{},
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

type GroupTable struct {
	ID                uuid.UUID            `gorm:"type:varchar(36);not null;primaryKey"`
	Name              string               `gorm:"type:varchar(64);size:64;not null"`
	GroupTypeID       int                  `gorm:"type:tinyint;not null"`
	Description       string               `gorm:"type:varchar(400);size:400;not null"`
	ReadPermissionID  int                  `gorm:"type:tinyint;not null"`
	WritePermissionID int                  `gorm:"type:tinyint;not null"`
	CreatedAt         time.Time            `gorm:"type:datetime;not null"`
	GroupType         GroupTypeTable       `gorm:"foreignKey:GroupTypeID"`
	Administrators    []AdministratorTable `gorm:"foreignKey:GroupID"`
}

func (gt *GroupTable) TableName() string {
	return "groups"
}

type GroupTypeTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (gtt *GroupTypeTable) TableName() string {
	return "group_types"
}

type ReadPermissionTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (rpt *ReadPermissionTable) TableName() string {
	return "read_permissions"
}

type WritePermissionTable struct {
	ID     int    `gorm:"type:TINYINT AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);size:32;not null;unique"`
	Active bool   `gorm:"type:boolean;default:true"`
}

func (wpt *WritePermissionTable) TableName() string {
	return "write_permissions"
}

type AdministratorTable struct {
	GroupID uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
	UserID  uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
}

func (at *AdministratorTable) TableName() string {
	return "administrators"
}
