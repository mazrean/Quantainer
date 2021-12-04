package service

import (
	"context"
	"io"

	"github.com/mazrean/Quantainer/domain"
)

type File interface {
	Upload(ctx context.Context, reader io.Reader) (*domain.File, error)
}
