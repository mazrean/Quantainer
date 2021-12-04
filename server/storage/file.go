package storage

import (
	"context"
	"io"

	"github.com/mazrean/Quantainer/domain"
)

type File interface {
	SaveFile(ctx context.Context, file *domain.File, reader io.Reader) error
	GetFile(ctx context.Context, fileID string, writer io.Writer) error
}
