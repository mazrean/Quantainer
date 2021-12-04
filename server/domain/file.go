package domain

import (
	"time"

	"github.com/mazrean/Quantainer/domain/values"
)

type File struct {
	id        values.FileID
	fileType  values.FileType
	createdAt time.Time
}

func NewFile(
	id values.FileID,
	fileType values.FileType,
	createdAt time.Time,
) *File {
	return &File{
		id:        id,
		fileType:  fileType,
		createdAt: createdAt,
	}
}

func (f *File) GetID() values.FileID {
	return f.id
}

func (f *File) GetFileType() values.FileType {
	return f.fileType
}

func (f *File) GetCreatedAt() time.Time {
	return f.createdAt
}
