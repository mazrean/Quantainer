package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
)

type File interface {
	SaveFile(ctx context.Context, file *domain.File) error
}
