package v1

import (
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/storage"
)

type Resource struct {
	dbRepository       repository.DB
	fileRepository     repository.File
	resourceRepository repository.Resource
	fileStorage        storage.File
}

func NewResource(
	dbRepository repository.DB,
	fileRepository repository.File,
	resourceRepository repository.Resource,
	fileStorage storage.File,
) *Resource {
	return &Resource{
		dbRepository:       dbRepository,
		fileRepository:     fileRepository,
		resourceRepository: resourceRepository,
		fileStorage:        fileStorage,
	}
}
