package service

import (
	"context"
	"io"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
)

type File interface {
	Upload(ctx context.Context, session *domain.OIDCSession, reader io.Reader) (*FileInfo, error)
	UploadBotFile(ctx context.Context, user *UserInfo, reader io.Reader) (*FileInfo, error)
	Download(ctx context.Context, fileID values.FileID, writer io.Writer) (*domain.File, error)
}

type FileInfo struct {
	File    *domain.File
	Creator *UserInfo
}
