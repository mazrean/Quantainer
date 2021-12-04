package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/service"
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

func (r *Resource) CreateResource(
	ctx context.Context,
	user *domain.TraPMember,
	fileID values.FileID,
	name values.ResourceName,
	resourceType values.ResourceType,
	comment values.ResourceComment,
) (*domain.Resource, error) {
	var resource *domain.Resource
	err := r.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		_, err := r.fileRepository.GetFile(ctx, fileID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoFile
		}
		if err != nil {
			return fmt.Errorf("failed to get file: %w", err)
		}

		resource = domain.NewResource(
			values.NewResourceID(),
			name,
			resourceType,
			comment,
			time.Now(),
		)

		err = r.resourceRepository.SaveResource(ctx, fileID, resource)
		if err != nil {
			return fmt.Errorf("failed to save resource: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (r *Resource) GetResource(ctx context.Context, resourceID values.ResourceID) (*domain.Resource, error) {
	resource, err := r.resourceRepository.GetResource(ctx, resourceID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoResource
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return resource, nil
}

func (r *Resource) GetResources(ctx context.Context, params *service.ResourceSearchParams) ([]*domain.Resource, error) {
	resources, err := r.resourceRepository.GetResources(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return resources, nil
}
