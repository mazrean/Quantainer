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
		session *domain.OIDCSession,
		fileID values.FileID,
		name values.ResourceName,
		resourceType values.ResourceType,
		comment values.ResourceComment,
	) (*ResourceInfo, error)
	GetResource(ctx context.Context, session *domain.OIDCSession, resourceID values.ResourceID) (*ResourceInfo, error)
	GetResources(ctx context.Context, session *domain.OIDCSession, params *ResourceSearchParams) ([]*ResourceInfo, error)
	DownloadResourceFile(ctx context.Context, resourceID values.ResourceID, writer io.Writer) (*domain.File, error)
}

type ResourceSearchParams struct {
	ResourceTypes []values.ResourceType
	Users         []values.TraPMemberName
	Limit         int
	Offset        int
}

type ResourceInfo struct {
	*domain.Resource
	Creator *UserInfo
}
