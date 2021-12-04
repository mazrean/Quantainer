package v1

import (
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
