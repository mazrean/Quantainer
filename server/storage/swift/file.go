package swift

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/storage"
)

type File struct {
	client *Client
}

func NewFile(client *Client) *File {
	return &File{
		client: client,
	}
}

func (gf *File) SaveFile(ctx context.Context, file *domain.File, reader io.Reader) error {
	fileKey := gf.fileKey(file)

	var contentType string
	switch file.GetType() {
	case values.FileTypeJpeg:
		contentType = "image/jpeg"
	case values.FileTypePng:
		contentType = "image/png"
	case values.FileTypeWebP:
		contentType = "image/webp"
	case values.FileTypeSvg:
		contentType = "image/svg+xml"
	case values.FileTypeGif:
		contentType = "image/gif"
	default:
		contentType = "application/octet-stream"
	}

	err := gf.client.saveFile(
		ctx,
		fileKey,
		contentType,
		"",
		reader,
	)
	if errors.Is(err, ErrAlreadyExists) {
		return storage.ErrAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (gf *File) GetFile(ctx context.Context, file *domain.File, writer io.Writer) error {
	fileKey := gf.fileKey(file)

	err := gf.client.loadFile(
		ctx,
		fileKey,
		writer,
	)
	if errors.Is(err, ErrNotFound) {
		return storage.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	return nil
}

func (gf *File) fileKey(file *domain.File) string {
	return fmt.Sprintf("files/%s", uuid.UUID(file.GetID()).String())
}
