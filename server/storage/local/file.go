package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/storage"
)

type File struct {
	fileRootPath     string
	directoryManager *DirectoryManager
}

func NewFile(directoryManager *DirectoryManager) (*File, error) {
	fileRootPath, err := directoryManager.setupDirectory("files")
	if err != nil {
		return nil, fmt.Errorf("failed to setup directory: %w", err)
	}

	return &File{
		fileRootPath:     fileRootPath,
		directoryManager: directoryManager,
	}, nil
}

func (f *File) SaveFile(ctx context.Context, file *domain.File, reader io.Reader) error {
	filePath := path.Join(f.fileRootPath, uuid.UUID(file.GetID()).String())

	_, err := os.Stat(filePath)
	if err == nil {
		return storage.ErrAlreadyExists
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	fl, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fl.Close()

	_, err = io.Copy(fl, reader)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	return nil
}

func (f *File) GetFile(ctx context.Context, file *domain.File, writer io.Writer) error {
	filePath := path.Join(f.fileRootPath, uuid.UUID(file.GetID()).String())

	fl, err := os.Open(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		return storage.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fl.Close()

	_, err = io.Copy(writer, fl)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	return nil
}
