package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
)

type File interface {
	SaveFile(ctx context.Context, file *domain.File) error
	GetFile(ctx context.Context, fileID values.FileID, lockType LockType) (*domain.File, error)
	GetFileByResourceID(ctx context.Context, resourceID values.ResourceID) (*domain.File, error)
}
