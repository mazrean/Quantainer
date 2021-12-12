package service

import (
	"context"
	"time"

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
	CreateBotResource(
		ctx context.Context,
		user *UserInfo,
		fileID values.FileID,
		name values.ResourceName,
		resourceType values.ResourceType,
		comment values.ResourceComment,
		createdAt time.Time,
	) (*ResourceInfo, error)
	GetResource(ctx context.Context, session *domain.OIDCSession, resourceID values.ResourceID) (*ResourceInfo, error)
	GetResources(ctx context.Context, session *domain.OIDCSession, params *ResourceSearchParams) ([]*ResourceInfo, error)
}

type ResourceSearchParams struct {
	ResourceTypes []values.ResourceType
	Users         []values.TraPMemberName
	Group         *values.GroupID
	Limit         int
	Offset        int
}

type ResourceInfo struct {
	*domain.Resource
	*domain.File
	Creator *UserInfo
}
