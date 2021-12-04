package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type Resource interface {
	SaveResource(ctx context.Context, fileID values.FileID, resource *domain.Resource) error
	GetResource(ctx context.Context, resourceID values.ResourceID) (*ResourceWithCreator, error)
	GetResources(ctx context.Context, params *service.ResourceSearchParams) ([]*ResourceWithCreator, error)
}

type ResourceWithCreator struct {
	*domain.Resource
	Creator values.TraPMemberID
}
