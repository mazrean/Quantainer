package gorm2

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Group struct {
	db *DB
}

func NewGroup(db *DB) (*Group, error) {
	ctx := context.Background()

	gormDB, err := db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	err = setupGroupTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup group type table: %w", err)
	}

	err = setupReadPermissionTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup read permission type table: %w", err)
	}

	err = setupWritePermissionTypeTable(gormDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup write permission type table: %w", err)
	}

	return &Group{
		db: db,
	}, nil
}

const (
	groupTypeArtBook = "art_book"
	groupTypeOther   = "other"
)

func setupGroupTypeTable(db *gorm.DB) error {
	groupTypes := []GroupTypeTable{
		{
			Name:   groupTypeArtBook,
			Active: true,
		},
		{
			Name:   groupTypeOther,
			Active: true,
		},
	}

	for _, groupType := range groupTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", groupType.Name).
			FirstOrCreate(&groupType).Error
		if err != nil {
			return fmt.Errorf("failed to create group type: %w", err)
		}
	}

	return nil
}

const (
	readPermissionPublic  = "public"
	readPermissionPrivate = "private"
)

func setupReadPermissionTypeTable(db *gorm.DB) error {
	readPermissions := []ReadPermissionTable{
		{
			Name:   readPermissionPublic,
			Active: true,
		},
		{
			Name:   readPermissionPrivate,
			Active: true,
		},
	}

	for _, readPermission := range readPermissions {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", readPermission.Name).
			FirstOrCreate(&readPermission).Error
		if err != nil {
			return fmt.Errorf("failed to create read permission type: %w", err)
		}
	}

	return nil
}

const (
	writePermissionPublic  = "public"
	writePermissionPrivate = "private"
)

func setupWritePermissionTypeTable(db *gorm.DB) error {
	writePermissions := []WritePermissionTable{
		{
			Name:   writePermissionPublic,
			Active: true,
		},
		{
			Name:   writePermissionPrivate,
			Active: true,
		},
	}

	for _, writePermission := range writePermissions {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", writePermission.Name).
			FirstOrCreate(&writePermission).Error
		if err != nil {
			return fmt.Errorf("failed to create read permission type: %w", err)
		}
	}

	return nil
}
