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
	"github.com/mazrean/Quantainer/service"
	"github.com/mazrean/Quantainer/storage"
)

type File struct {
	dbRepository   repository.DB
	fileRepository repository.File
	fileStorage    storage.File
	userUtils      *UserUtils
}

func NewFile(
	dbRepository repository.DB,
	fileRepository repository.File,
	fileStorage storage.File,
	userUtils *UserUtils,
) *File {
	return &File{
		dbRepository:   dbRepository,
		fileRepository: fileRepository,
		fileStorage:    fileStorage,
		userUtils:      userUtils,
	}
}

func (f *File) Upload(ctx context.Context, session *domain.OIDCSession, reader io.Reader) (*service.FileInfo, error) {
	user, err := f.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	tr := io.TeeReader(reader, buf)

	content, err := io.ReadAll(tr)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
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

	file := domain.NewFile(
		values.NewFileID(),
		fileType,
		time.Now(),
	)

	err = f.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		err := f.fileRepository.SaveFile(ctx, user, file)
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

	return &service.FileInfo{
		File:    file,
		Creator: user,
	}, nil
}
