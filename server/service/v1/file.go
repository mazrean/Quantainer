package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	svg "github.com/h2non/go-is-svg"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/storage"
)

type File struct {
	dbRepository   repository.DB
	fileRepository repository.File
	fileStorage    storage.File
}

func NewFile(
	dbRepository repository.DB,
	fileRepository repository.File,
	fileStorage storage.File,
) *File {
	return &File{
		dbRepository:   dbRepository,
		fileRepository: fileRepository,
		fileStorage:    fileStorage,
	}
}

func (f *File) Upload(ctx context.Context, reader io.Reader) (*domain.File, error) {
	var file *domain.File
	err := f.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		buf := bytes.NewBuffer(nil)
		tr := io.TeeReader(reader, buf)

		content, err := io.ReadAll(tr)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var fileType values.FileType
		mime := http.DetectContentType(content)
		switch mime {
		case "image/jpeg":
			fileType = values.FileTypeJpeg
		case "image/png":
			fileType = values.FileTypePng
		case "image/webp":
			fileType = values.FileTypeWebP
		case "image/gif":
			fileType = values.FileTypeGif
		default:
			if svg.Is(content) {
				fileType = values.FileTypeSvg
			} else {
				fileType = values.FileTypeOther
			}
		}

		file = domain.NewFile(
			values.NewFileID(),
			fileType,
			time.Now(),
		)

		err = f.fileRepository.SaveFile(ctx, file)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		err = f.fileStorage.SaveFile(ctx, file, buf)
		if err != nil {
			return fmt.Errorf("failed to save file: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return file, nil
}
