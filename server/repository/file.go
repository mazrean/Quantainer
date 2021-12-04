package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type File interface {
	SaveFile(ctx context.Context, user *service.UserInfo, file *domain.File) error
	GetFile(ctx context.Context, fileID values.FileID, lockType LockType) (*domain.File, error)
	GetFileByResourceID(ctx context.Context, resourceID values.ResourceID) (*domain.File, error)
}
