package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type File interface {
	SaveFile(ctx context.Context, user *service.UserInfo, file *domain.File) error
	GetFile(ctx context.Context, fileID values.FileID, lockType LockType) (*FileWithCreator, error)
}

type FileWithCreator struct {
	*domain.File
	Creator values.TraPMemberID
}
