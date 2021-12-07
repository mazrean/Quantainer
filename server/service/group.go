package service

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
)

type Group interface {
	CreateGroup(
		ctx context.Context,
		session *domain.OIDCSession,
		name values.GroupName,
		groupType values.GroupType,
		description values.GroupDescription,
		readPermission values.GroupReadPermission,
		writePermission values.GroupWritePermission,
		mainResource values.ResourceID,
		resources []values.ResourceID,
	) (*domain.Group, *ResourceInfo, error)
	EditGroup(
		ctx context.Context,
		session *domain.OIDCSession,
		id values.GroupID,
		name values.GroupName,
		groupType values.GroupType,
		description values.GroupDescription,
		readPermission values.GroupReadPermission,
		writePermission values.GroupWritePermission,
		mainResource values.ResourceID,
		resources []values.ResourceID,
	) (*GroupDetail, error)
	DeleteGroup(ctx context.Context, session *domain.OIDCSession, id values.GroupID) error
	AddResource(ctx context.Context, session *domain.OIDCSession, id values.GroupID, resource values.ResourceID) ([]*ResourceInfo, error)
	GetGroup(ctx context.Context, session *domain.OIDCSession, groupID values.GroupID) (*GroupDetail, error)
	GetGroups(ctx context.Context, session *domain.OIDCSession, params *GroupSearchParams) ([]*GroupInfo, error)
}

type GroupInfo struct {
	*domain.Group
	MainResource *ResourceInfo
}

type GroupDetail struct {
	*domain.Group
	Administers  []*UserInfo
	MainResource *ResourceInfo
}

type GroupSearchParams struct {
	GroupTypes []values.GroupType
	Users      []values.TraPMemberName
	Limit      int
	Offset     int
}
