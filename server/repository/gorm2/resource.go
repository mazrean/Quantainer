package gorm2

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	resourceTypeImage = "image"
	resourceTypeOther = "other"
)

type Resource struct {
	db *DB
}

func NewResource(db *DB) (*Resource, error) {
	ctx := context.Background()

	gormDB, err := db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	err = setupResourceTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup resource type table: %w", err)
	}

	return &Resource{
		db: db,
	}, nil
}

func setupResourceTypeTable(db *gorm.DB) error {
	resourceTypes := []ResourceTypeTable{
		{
			Name:   resourceTypeImage,
			Active: true,
		},
		{
			Name:   resourceTypeOther,
			Active: true,
		},
	}

	for _, resourceType := range resourceTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", resourceType.Name).
			FirstOrCreate(&resourceType).Error
		if err != nil {
			return fmt.Errorf("failed to create resource type: %w", err)
		}
	}

	return nil
}
