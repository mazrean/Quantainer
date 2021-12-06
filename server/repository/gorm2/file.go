package gorm2

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/service"
	"gorm.io/gorm"
)

const (
	fileTypeJpeg  = "jpeg"
	fileTypePng   = "png"
	fileTypeWebP  = "webp"
	fileTypeSvg   = "svg"
	fileTypeGif   = "gif"
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

func (f *File) SaveFile(ctx context.Context, user *service.UserInfo, file *domain.File) error {
	db, err := f.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	var fileTypeName string
	switch file.GetType() {
	case values.FileTypeJpeg:
		fileTypeName = fileTypeJpeg
	case values.FileTypePng:
		fileTypeName = fileTypePng
	case values.FileTypeWebP:
		fileTypeName = fileTypeWebP
	case values.FileTypeSvg:
		fileTypeName = fileTypeSvg
	case values.FileTypeGif:
		fileTypeName = fileTypeGif
	case values.FileTypeOther:
		fileTypeName = fileTypeOther
	default:
		return fmt.Errorf("invalid file type: %d", file.GetType())
	}

	var fileType FileTypeTable
	err = db.
		Session(&gorm.Session{}).
		Where("name = ?", fileTypeName).
		Where("active").
		Select("id").
		Take(&fileType).Error
	if err != nil {
		return fmt.Errorf("failed to get file type: %w", err)
	}
	fileTypeID := fileType.ID

	fileTable := FileTable{
		ID:         uuid.UUID(file.GetID()),
		FileTypeID: fileTypeID,
		CreatorID:  uuid.UUID(user.GetID()),
		CreatedAt:  file.GetCreatedAt(),
	}

	err = db.Create(&fileTable).Error
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (f *File) GetFile(ctx context.Context, fileID values.FileID, lockType repository.LockType) (*repository.FileWithCreator, error) {
	db, err := f.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	db, err = f.db.setLock(db, lockType)
	if err != nil {
		return nil, fmt.Errorf("failed to set lock: %w", err)
	}

	var fileTable FileTable
	err = db.
		Session(&gorm.Session{}).
		Joins("FileType").
		Where("files.id = ?", fileID).
		Select("files.creator_id", "files.created_at", "file_types.name").
		Take(&fileTable).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file type: %w", err)
	}

	var fileType values.FileType
	switch fileTable.FileType.Name {
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
		return nil, fmt.Errorf("invalid file type: %s", fileTable.FileType.Name)
	}

	return &repository.FileWithCreator{
		File: domain.NewFile(
			fileID,
			fileType,
			fileTable.CreatedAt,
		),
		Creator: values.NewTrapMemberID(fileTable.CreatorID),
	}, nil
}
