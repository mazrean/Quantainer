package gorm2

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	fileTypeJpeg = "jpeg"
	fileTypePng  = "png"
	fileTypeWebP = "webp"
	fileTypeSvg = "svg"
	fileTypeGif  = "gif"
	fileTypeOther = "other"
)

type File struct {
	db *DB
}

func NewFile(db *DB) (*File, error) {
	ctx := context.Background()

	gormDB, err := db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	err = setupFileTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup file type table: %w", err)
	}

	return &File{
		db: db,
	}, nil
}

func setupFileTypeTable(db *gorm.DB) error {
	fileTypes := []FileTypeTable{
		{
			Name:   fileTypeJpeg,
			Active: true,
		},
		{
			Name:   fileTypePng,
			Active: true,
		},
		{
			Name:   fileTypeWebP,
			Active: true,
		},
		{
			Name:   fileTypeSvg,
			Active: true,
		},
		{
			Name:   fileTypeGif,
			Active: true,
		},
		{
			Name:   fileTypeOther,
			Active: true,
		},
	}

	for _, fileType := range fileTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", fileType.Name).
			FirstOrCreate(&fileType).Error
		if err != nil {
			return fmt.Errorf("failed to create role type: %w", err)
		}
	}

	return nil
}
