package values

import (
	"github.com/google/uuid"
)

type (
	FileID   uuid.UUID
	FileType int8
)

func NewFileID() FileID {
	return FileID(uuid.New())
}

func NewFileIDFromUUID(u uuid.UUID) FileID {
	return FileID(u)
}

const (
	FileTypeJpeg FileType = iota + 1
	FileTypePng
	FileTypeWebP
	FileTypeSvg
	FileTypeGif
	FileTypeOther
)

func (ft FileType) IsValidResourceType(resourceType ResourceType) bool {
	switch ft {
	case FileTypeJpeg, FileTypePng, FileTypeWebP, FileTypeSvg, FileTypeGif:
		return resourceType == ResourceTypeImage || resourceType == ResourceTypeOther
	default:
		return resourceType == ResourceTypeOther
	}
}
