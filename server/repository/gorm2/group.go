package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
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

func (g *Group) SaveGroup(ctx context.Context, group *domain.Group, mainResource values.ResourceID) error {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	var groupTypeName string
	switch group.GetType() {
	case values.GroupTypeArtBook:
		groupTypeName = groupTypeArtBook
	case values.GroupTypeOther:
		groupTypeName = groupTypeOther
	default:
		return fmt.Errorf("invalid group type: %d", group.GetType())
	}

	var groupType GroupTypeTable
	err = db.
		Session(&gorm.Session{}).
		Where("name = ?", groupTypeName).
		Take(&groupType).Error
	if err != nil {
		return fmt.Errorf("failed to get group type: %w", err)
	}

	var readPermissionName string
	switch group.GetReadPermission() {
	case values.GroupReadPermissionPublic:
		readPermissionName = readPermissionPublic
	case values.GroupReadPermissionPrivate:
		readPermissionName = readPermissionPrivate
	default:
		return fmt.Errorf("invalid read permission: %d", group.GetReadPermission())
	}

	var readPermission ReadPermissionTable
	err = db.
		Session(&gorm.Session{}).
		Where("name = ?", readPermissionName).
		Take(&readPermission).Error
	if err != nil {
		return fmt.Errorf("failed to get read permission: %w", err)
	}

	var writePermissionName string
	switch group.GetWritePermission() {
	case values.GroupWritePermissionPublic:
		writePermissionName = writePermissionPublic
	case values.GroupWritePermissionPrivate:
		writePermissionName = writePermissionPrivate
	default:
		return fmt.Errorf("invalid write permission: %d", group.GetWritePermission())
	}

	var writePermission WritePermissionTable
	err = db.
		Session(&gorm.Session{}).
		Where("name = ?", writePermissionName).
		Take(&writePermission).Error
	if err != nil {
		return fmt.Errorf("failed to get write permission: %w", err)
	}

	groupTable := GroupTable{
		ID:                uuid.UUID(group.GetID()),
		Name:              string(group.GetName()),
		Description:       string(group.GetDescription()),
		GroupTypeID:       groupType.ID,
		MainResourceID:    uuid.UUID(mainResource),
		ReadPermissionID:  readPermission.ID,
		WritePermissionID: writePermission.ID,
	}

	err = db.Create(&groupTable).Error
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}
