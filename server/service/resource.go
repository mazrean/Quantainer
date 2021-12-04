package service

import (
	"context"
	"io"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
)

type Resource interface {
	CreateResource(
		ctx context.Context,
		user *domain.TraPMember,
		fileID values.FileID,
		name values.ResourceName,
		resourceType values.ResourceType,
		comment values.ResourceComment,
	) (*domain.Resource, error)
	GetResource(ctx context.Context, resourceID values.ResourceID) (*domain.Resource, error)
	GetResources(ctx context.Context, params *ResourceSearchParams) ([]*domain.Resource, error)
	DownloadResourceFile(ctx context.Context, resourceID values.ResourceID, writer io.Writer) (*domain.File, error)
}

type ResourceSearchParams struct {
	ResourceTypes []values.ResourceType
	Users         []values.TraPMemberName
	Limit         int
	Offset        int
}
