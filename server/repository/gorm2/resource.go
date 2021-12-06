package gorm2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
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

func (r *Resource) SaveResource(ctx context.Context, fileID values.FileID, resource *domain.Resource) error {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	var resourceTypeName string
	switch resource.GetType() {
	case values.ResourceTypeImage:
		resourceTypeName = resourceTypeImage
	case values.ResourceTypeOther:
		resourceTypeName = resourceTypeOther
	default:
		return fmt.Errorf("invalid resource type: %d", resource.GetType())
	}

	resourceType := ResourceTypeTable{}
	err = db.
		Session(&gorm.Session{}).
		Where("name = ?", resourceTypeName).
		First(&resourceType).Error
	if err != nil {
		return fmt.Errorf("failed to get resource type: %w", err)
	}

	resourceTable := ResourceTable{
		ID:             uuid.UUID(resource.GetID()),
		FileID:         uuid.UUID(fileID),
		Name:           string(resource.GetName()),
		ResourceTypeID: resourceType.ID,
		Comment:        string(resource.GetComment()),
		CreatedAt:      resource.GetCreatedAt(),
	}

	err = db.Create(&resourceTable).Error
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return nil
}

func (r *Resource) GetResource(ctx context.Context, resourceID values.ResourceID) (*repository.ResourceWithCreator, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var resourceTable ResourceTable
	err = db.
		Session(&gorm.Session{}).
		Joins("ResourceType").
		Joins("File").
		Where("id = ?", uuid.UUID(resourceID)).
		Select("resources.name", "resources.comment", "resources.created_at", "resource_types.name", "files.creator_id").
		Take(&resourceTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	var resourceType values.ResourceType
	switch resourceTable.ResourceType.Name {
	case resourceTypeImage:
		resourceType = values.ResourceTypeImage
	case resourceTypeOther:
		resourceType = values.ResourceTypeOther
	default:
		return nil, fmt.Errorf("invalid resource type: %s", resourceTable.ResourceType.Name)
	}

	resource := repository.ResourceWithCreator{
		Resource: domain.NewResource(
			resourceID,
			values.NewResourceName(resourceTable.Name),
			resourceType,
			values.NewResourceComment(resourceTable.Comment),
			resourceTable.CreatedAt,
		),
		Creator: values.NewTrapMemberID(resourceTable.File.CreatorID),
	}

	return &resource, nil
}
