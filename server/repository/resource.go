package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type Resource interface {
	SaveResource(ctx context.Context, fileID values.FileID, resource *domain.Resource) error
	GetResource(ctx context.Context, resourceID values.ResourceID) (*ResourceInfo, error)
	GetResources(ctx context.Context, params *ResourceSearchParams) ([]*ResourceInfo, error)
	GetResourcesByIDs(ctx context.Context, resourceIDs []values.ResourceID, lockType LockType) ([]*domain.Resource, error)
}

type ResourceInfo struct {
	*domain.Resource
	*domain.File
	Creator values.TraPMemberID
}

type ResourceSearchParams struct {
	ResourceTypes []values.ResourceType
	Users         []*service.UserInfo
	Groups        []*domain.Group
	Limit         int
	Offset        int
}
