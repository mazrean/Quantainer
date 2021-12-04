package service

import (
	"context"
	"io"

	"github.com/mazrean/Quantainer/domain"
)

type File interface {
	Upload(ctx context.Context, session *domain.OIDCSession, reader io.Reader) (*FileInfo, error)
}

type FileInfo struct {
	File    *domain.File
	Creator *UserInfo
}
