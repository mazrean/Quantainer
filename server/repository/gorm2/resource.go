package gorm2

import (
	"context"
	"errors"
	"fmt"
	"log"

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
		Take(&resourceType).Error
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

func (r *Resource) GetResource(ctx context.Context, resourceID values.ResourceID) (*repository.ResourceInfo, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var resourceTable ResourceTable
	err = db.
		Session(&gorm.Session{}).
		Joins("ResourceType").
		Joins("File").
		Where("resources.id = ?", uuid.UUID(resourceID)).
		Select(
			"resources.name",
			"resources.comment",
			"resources.created_at",
		).
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

	err = db.
		Session(&gorm.Session{}).
		Where("id = ?", resourceTable.File.FileTypeID).
		Take(&resourceTable.File.FileType).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get resource type: %w", err)
	}

	var fileType values.FileType
	switch resourceTable.File.FileType.Name {
	case fileTypeJpeg:
		fileType = values.FileTypeJpeg
	case fileTypePng:
		fileType = values.FileTypePng
	case fileTypeWebP:
		fileType = values.FileTypeWebP
	case fileTypeSvg:
		fileType = values.FileTypeSvg
	case fileTypeGif:
		fileType = values.FileTypeGif
	case fileTypeOther:
		fileType = values.FileTypeOther
	default:
		log.Printf("debug: resource: %+v\n", resourceTable)
		return nil, fmt.Errorf("invalid file type: %s", resourceTable.File.FileType.Name)
	}

	resource := repository.ResourceInfo{
		Resource: domain.NewResource(
			resourceID,
			values.NewResourceName(resourceTable.Name),
			resourceType,
			values.NewResourceComment(resourceTable.Comment),
			resourceTable.CreatedAt,
		),
		File: domain.NewFile(
			values.NewFileIDFromUUID(resourceTable.File.ID),
			fileType,
			resourceTable.File.CreatedAt,
		),
		Creator: values.NewTrapMemberID(resourceTable.File.CreatorID),
	}

	return &resource, nil
}

func (r *Resource) GetResources(ctx context.Context, params *repository.ResourceSearchParams) ([]*repository.ResourceInfo, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	resourceTypeNames := make([]string, 0, len(params.ResourceTypes))
	for _, resourceType := range params.ResourceTypes {
		switch resourceType {
		case values.ResourceTypeImage:
			resourceTypeNames = append(resourceTypeNames, resourceTypeImage)
		case values.ResourceTypeOther:
			resourceTypeNames = append(resourceTypeNames, resourceTypeOther)
		default:
			return nil, fmt.Errorf("invalid resource type: %d", resourceType)
		}
	}

	creatorIDs := make([]uuid.UUID, 0, len(params.Users))
	for _, creatorInfo := range params.Users {
		creatorIDs = append(creatorIDs, uuid.UUID(creatorInfo.GetID()))
	}

	query := db.
		Session(&gorm.Session{}).
		Joins("ResourceType").
		Joins("File")

	if len(resourceTypeNames) != 0 {
		query = query.Where("ResourceType.name IN ?", resourceTypeNames)
	}
	if len(creatorIDs) != 0 {
		query = query.Where("File.creator_id IN ?", creatorIDs)
	}

	if len(params.Groups) != 0 {
		groupIDs := make([]uuid.UUID, 0, len(params.Groups))
		for _, groupInfo := range params.Groups {
			groupIDs = append(groupIDs, uuid.UUID(groupInfo.GetID()))
		}
		query = query.
			Joins("JOIN group_resources ON resources.id = group_resources.resource_table_id").
			Where("group_resources.id IN ?", groupIDs)
	}

	if params.Limit != -1 {
		query = query.Limit(params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	var resourceTables []ResourceTable
	err = query.
		Select(
			"resources.id",
			"resources.name",
			"resources.comment",
			"resources.created_at",
		).
		Find(&resourceTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	var fileTypeTables []FileTypeTable
	err = db.
		Session(&gorm.Session{}).
		Find(&fileTypeTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get resource type: %w", err)
	}

	fileTypeMap := make(map[int]string, len(fileTypeTables))
	for _, fileTypeTable := range fileTypeTables {
		fileTypeMap[fileTypeTable.ID] = fileTypeTable.Name
	}

	resources := make([]*repository.ResourceInfo, 0, len(resourceTables))
	for _, resourceTable := range resourceTables {
		var resourceType values.ResourceType
		switch resourceTable.ResourceType.Name {
		case resourceTypeImage:
			resourceType = values.ResourceTypeImage
		case resourceTypeOther:
			resourceType = values.ResourceTypeOther
		default:
			return nil, fmt.Errorf("invalid resource type: %s", resourceTable.ResourceType.Name)
		}

		var fileType values.FileType
		switch fileTypeMap[resourceTable.File.FileTypeID] {
		case fileTypeJpeg:
			fileType = values.FileTypeJpeg
		case fileTypePng:
			fileType = values.FileTypePng
		case fileTypeWebP:
			fileType = values.FileTypeWebP
		case fileTypeSvg:
			fileType = values.FileTypeSvg
		case fileTypeGif:
			fileType = values.FileTypeGif
		case fileTypeOther:
			fileType = values.FileTypeOther
		default:
			return nil, fmt.Errorf("invalid file type: %s", resourceTable.File.FileType.Name)
		}

		resource := repository.ResourceInfo{
			Resource: domain.NewResource(
				values.ResourceID(resourceTable.ID),
				values.NewResourceName(resourceTable.Name),
				resourceType,
				values.NewResourceComment(resourceTable.Comment),
				resourceTable.CreatedAt,
			),
			File: domain.NewFile(
				values.NewFileIDFromUUID(resourceTable.File.ID),
				fileType,
				resourceTable.File.CreatedAt,
			),
			Creator: values.NewTrapMemberID(resourceTable.File.CreatorID),
		}

		resources = append(resources, &resource)
	}

	return resources, nil
}

func (r *Resource) GetResourcesByIDs(ctx context.Context, resourceIDs []values.ResourceID, lockType repository.LockType) ([]*domain.Resource, error) {
	db, err := r.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	uuidResourceIDs := make([]uuid.UUID, 0, len(resourceIDs))
	for _, resourceID := range resourceIDs {
		uuidResourceIDs = append(uuidResourceIDs, uuid.UUID(resourceID))
	}

	var resourceTables []ResourceTable
	err = db.
		Session(&gorm.Session{}).
		Joins("ResourceType").
		Where("resources.id IN (?)", uuidResourceIDs).
		Find(&resourceTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	resources := make([]*domain.Resource, 0, len(resourceTables))
	for _, resourceTable := range resourceTables {
		var resourceType values.ResourceType
		switch resourceTable.ResourceType.Name {
		case resourceTypeImage:
			resourceType = values.ResourceTypeImage
		case resourceTypeOther:
			resourceType = values.ResourceTypeOther
		default:
			return nil, fmt.Errorf("invalid resource type: %s", resourceTable.ResourceType.Name)
		}

		resources = append(resources, domain.NewResource(
			values.NewResourceIDFromUUID(resourceTable.ID),
			values.NewResourceName(resourceTable.Name),
			resourceType,
			values.NewResourceComment(resourceTable.Comment),
			resourceTable.CreatedAt,
		))
	}

	return resources, nil
}
