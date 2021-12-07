package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type Group interface {
	SaveGroup(ctx context.Context, group *domain.Group, mainResource values.ResourceID) error
	EditGroup(ctx context.Context, group *domain.Group, mainResource values.ResourceID) error
	DeleteGroup(ctx context.Context, group *domain.Group) error
	AddResources(ctx context.Context, group *domain.Group, resources []values.ResourceID) error
	DeleteResources(ctx context.Context, group *domain.Group, resources []values.ResourceID) error
	GetGroup(ctx context.Context, groupID values.GroupID, lockType LockType) (*GroupInfo, error)
	GetGroups(ctx context.Context, user *service.UserInfo, params *GroupSearchParams) ([]*GroupInfo, error)
}

type GroupInfo struct {
	*domain.Group
	MainResource *ResourceInfo
}

type GroupSearchParams struct {
	GroupTypes []values.GroupType
	Users      []*service.UserInfo
	Limit      int
	Offset     int
}
