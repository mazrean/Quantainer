package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/service"
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

func (g *Group) EditGroup(ctx context.Context, group *domain.Group, mainResource values.ResourceID) error {
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

	result := db.Updates(&groupTable)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	if result.RowsAffected == 0 {
		return repository.ErrNoRecordUpdated
	}

	return nil
}

func (g *Group) DeleteGroup(ctx context.Context, group *domain.Group) error {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	err = db.
		Where("id = ?", uuid.UUID(group.GetID())).
		Delete(&GroupTable{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

func (g *Group) AddResources(ctx context.Context, group *domain.Group, resources []values.ResourceID) error {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	groupTable := GroupTable{
		ID: uuid.UUID(group.GetID()),
	}

	resourceTables := make([]*ResourceTable, 0, len(resources))
	for _, resource := range resources {
		resourceTables = append(resourceTables, &ResourceTable{
			ID: uuid.UUID(resource),
		})
	}

	err = db.
		Model(&groupTable).
		Association("Resources").
		Append(resourceTables)
	if err != nil {
		return fmt.Errorf("failed to add resources to group: %w", err)
	}

	return nil
}

func (g *Group) DeleteResources(ctx context.Context, group *domain.Group, resources []values.ResourceID) error {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	groupTable := GroupTable{
		ID: uuid.UUID(group.GetID()),
	}

	resourceTables := make([]*ResourceTable, 0, len(resources))
	for _, resource := range resources {
		resourceTables = append(resourceTables, &ResourceTable{
			ID: uuid.UUID(resource),
		})
	}

	err = db.
		Model(&groupTable).
		Association("Resources").
		Delete(resourceTables)
	if err != nil {
		return fmt.Errorf("failed to delete resources from group: %w", err)
	}

	return nil
}

func (g *Group) GetGroup(ctx context.Context, groupID values.GroupID, lockType repository.LockType) (*repository.GroupInfo, error) {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = g.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var groupTable GroupTable
	err = db.
		Joins("GroupType").
		Joins("ReadPermission").
		Joins("WritePermission").
		Joins("MainResource").
		Where("id = ?", uuid.UUID(groupID)).
		Take(&groupTable).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	var groupType values.GroupType
	switch groupTable.GroupType.Name {
	case groupTypeArtBook:
		groupType = values.GroupTypeArtBook
	case groupTypeOther:
		groupType = values.GroupTypeOther
	default:
		return nil, fmt.Errorf("invalid group type: %s", groupTable.GroupType.Name)
	}

	var readPermission values.GroupReadPermission
	switch groupTable.ReadPermission.Name {
	case readPermissionPublic:
		readPermission = values.GroupReadPermissionPublic
	case readPermissionPrivate:
		readPermission = values.GroupReadPermissionPrivate
	default:
		return nil, fmt.Errorf("invalid read permission: %s", groupTable.ReadPermission.Name)
	}

	var writePermission values.GroupWritePermission
	switch groupTable.WritePermission.Name {
	case writePermissionPublic:
		writePermission = values.GroupWritePermissionPublic
	case writePermissionPrivate:
		writePermission = values.GroupWritePermissionPrivate
	default:
		return nil, fmt.Errorf("invalid write permission: %s", groupTable.WritePermission.Name)
	}

	var resourceTypeTable ResourceTypeTable
	err = db.
		Where("id = ?", groupTable.MainResource.ResourceTypeID).
		Take(&resourceTypeTable).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get main resource type: %w", err)
	}

	var resourceFileTable FileTable
	err = db.
		Joins("FileType").
		Where("id = ?", groupTable.MainResource.FileID).
		Take(&resourceFileTable).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get main resource file: %w", err)
	}

	var resourceType values.ResourceType
	switch resourceTypeTable.Name {
	case resourceTypeImage:
		resourceType = values.ResourceTypeImage
	case resourceTypeOther:
		resourceType = values.ResourceTypeOther
	default:
		return nil, fmt.Errorf("invalid resource type: %s", resourceTypeTable.Name)
	}

	var fileType values.FileType
	switch resourceFileTable.FileType.Name {
	case fileTypeJpeg:
		fileType = values.FileTypeJpeg
	case fileTypePng:
		fileType = values.FileTypePng
	case fileTypeOther:
		fileType = values.FileTypeOther
	default:
		return nil, fmt.Errorf("invalid file type: %s", resourceFileTable.FileType.Name)
	}

	return &repository.GroupInfo{
		Group: domain.NewGroup(
			values.NewGroupIDFromUUID(groupTable.ID),
			values.NewGroupName(groupTable.Name),
			groupType,
			values.NewGroupDescription(groupTable.Description),
			readPermission,
			writePermission,
			groupTable.CreatedAt,
		),
		MainResource: &repository.ResourceInfo{
			Resource: domain.NewResource(
				values.NewResourceIDFromUUID(groupTable.MainResource.ID),
				values.NewResourceName(groupTable.MainResource.Name),
				resourceType,
				values.NewResourceComment(groupTable.MainResource.Comment),
				groupTable.MainResource.CreatedAt,
			),
			File: domain.NewFile(
				values.NewFileIDFromUUID(resourceFileTable.ID),
				fileType,
				resourceFileTable.CreatedAt,
			),
			Creator: values.NewTrapMemberID(resourceFileTable.CreatorID),
		},
	}, nil
}

func (g *Group) GetGroups(ctx context.Context, user *service.UserInfo, params *repository.GroupSearchParams) ([]*repository.GroupInfo, error) {
	db, err := g.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	query := db.
		Joins("GroupType").
		Joins("ReadPermission").
		Joins("WritePermission").
		Preload("MainResource").
		Preload("MainResource.ResourceType").
		Preload("MainResource.File").
		Preload("MainResource.File.FileType")

	if len(params.GroupTypes) != 0 {
		groupTypeNames := make([]string, 0, len(params.GroupTypes))
		for _, groupType := range params.GroupTypes {
			switch groupType {
			case values.GroupTypeArtBook:
				groupTypeNames = append(groupTypeNames, groupTypeArtBook)
			case values.GroupTypeOther:
				groupTypeNames = append(groupTypeNames, groupTypeOther)
			default:
				return nil, fmt.Errorf("invalid group type: %d", groupType)
			}
		}

		query = query.Where("GroupType.name IN (?)", groupTypeNames)
	}

	if len(params.Users) != 0 {
		users := make([]uuid.UUID, 0, len(params.Users))
		for _, user := range params.Users {
			users = append(users, uuid.UUID(user.GetID()))
		}

		query = query.Where("MainResource.creator_id IN (?)", users)
	}

	if params.Limit != -1 {
		query = query.Limit(params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	var groupTables []GroupTable
	err = query.Find(&groupTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	var groups []*repository.GroupInfo
	for _, groupTable := range groupTables {
		var groupType values.GroupType
		switch groupTable.GroupType.Name {
		case groupTypeArtBook:
			groupType = values.GroupTypeArtBook
		case groupTypeOther:
			groupType = values.GroupTypeOther
		default:
			return nil, fmt.Errorf("invalid group type: %s", groupTable.GroupType.Name)
		}

		var readPermission values.GroupReadPermission
		switch groupTable.ReadPermission.Name {
		case readPermissionPublic:
			readPermission = values.GroupReadPermissionPublic
		case readPermissionPrivate:
			readPermission = values.GroupReadPermissionPrivate
		default:
			return nil, fmt.Errorf("invalid read permission: %s", groupTable.ReadPermission.Name)
		}

		var writePermission values.GroupWritePermission
		switch groupTable.WritePermission.Name {
		case writePermissionPublic:
			writePermission = values.GroupWritePermissionPublic
		case writePermissionPrivate:
			writePermission = values.GroupWritePermissionPrivate
		default:
			return nil, fmt.Errorf("invalid write permission: %s", groupTable.WritePermission.Name)
		}

		var resourceType values.ResourceType
		switch groupTable.MainResource.ResourceType.Name {
		case resourceTypeImage:
			resourceType = values.ResourceTypeImage
		case resourceTypeOther:
			resourceType = values.ResourceTypeOther
		default:
			return nil, fmt.Errorf("invalid resource type: %s", groupTable.MainResource.ResourceType.Name)
		}

		var fileType values.FileType
		switch groupTable.MainResource.File.FileType.Name {
		case fileTypeJpeg:
			fileType = values.FileTypeJpeg
		case fileTypePng:
			fileType = values.FileTypePng
		case fileTypeOther:
			fileType = values.FileTypeOther
		default:
			return nil, fmt.Errorf("invalid file type: %s", groupTable.MainResource.File.FileType.Name)
		}

		groups = append(groups, &repository.GroupInfo{
			Group: domain.NewGroup(
				values.NewGroupIDFromUUID(groupTable.ID),
				values.NewGroupName(groupTable.Name),
				groupType,
				values.NewGroupDescription(groupTable.Description),
				readPermission,
				writePermission,
				groupTable.CreatedAt,
			),
			MainResource: &repository.ResourceInfo{
				Resource: domain.NewResource(
					values.NewResourceIDFromUUID(groupTable.MainResource.ID),
					values.NewResourceName(groupTable.MainResource.Name),
					resourceType,
					values.NewResourceComment(groupTable.MainResource.Comment),
					groupTable.MainResource.CreatedAt,
				),
				File: domain.NewFile(
					values.NewFileIDFromUUID(groupTable.MainResource.File.ID),
					fileType,
					groupTable.MainResource.File.CreatedAt,
				),
				Creator: values.NewTrapMemberID(groupTable.MainResource.File.CreatorID),
			},
		})
	}

	return groups, nil
}
