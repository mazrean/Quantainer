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
)

type Group struct {
	dbRepository            repository.DB
	resourceRepository      repository.Resource
	groupRepository         repository.Group
	administratorRepository repository.Administrator
	userUtils               *UserUtils
}

func NewGroup(
	dbRepository repository.DB,
	resourceRepository repository.Resource,
	groupRepository repository.Group,
	administratorRepository repository.Administrator,
	userUtils *UserUtils,
) *Group {
	return &Group{
		dbRepository:            dbRepository,
		resourceRepository:      resourceRepository,
		groupRepository:         groupRepository,
		administratorRepository: administratorRepository,
		userUtils:               userUtils,
	}
}

func (g *Group) CreateGroup(
	ctx context.Context,
	session *domain.OIDCSession,
	name values.GroupName,
	groupType values.GroupType,
	description values.GroupDescription,
	readPermission values.GroupReadPermission,
	writePermission values.GroupWritePermission,
	mainResource values.ResourceID,
	resources []values.ResourceID,
) (*service.GroupDetail, error) {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	group := domain.NewGroup(
		values.NewGroupID(),
		name,
		groupType,
		description,
		readPermission,
		writePermission,
		time.Now(),
	)

	if !group.IsValidPermission() {
		return nil, service.ErrInvalidPermission
	}

	users, err := g.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var mainResourceInfo *service.ResourceInfo
	err = g.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		resourceInfo, err := g.resourceRepository.GetResource(ctx, mainResource)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoResource
		}
		if err != nil {
			return fmt.Errorf("failed to get main resource: %w", err)
		}

		var creator *service.UserInfo
		for i, user := range users {
			if user.GetID() == resourceInfo.Creator {
				creator = user
				break
			}

			if i == len(users)-1 {
				return service.ErrNoUser
			}
		}

		mainResourceInfo = &service.ResourceInfo{
			Resource: resourceInfo.Resource,
			File:     resourceInfo.File,
			Creator:  creator,
		}

		resourceList, err := g.resourceRepository.GetResourcesByIDs(ctx, resources, repository.LockTypeNone)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		if len(resourceList) != len(resources) {
			return service.ErrNoResource
		}

		err = g.groupRepository.SaveGroup(ctx, group, mainResource)
		if err != nil {
			return fmt.Errorf("failed to save group: %w", err)
		}

		err = g.groupRepository.AddResources(ctx, group, resources)
		if err != nil {
			return fmt.Errorf("failed to add resources: %w", err)
		}

		err = g.administratorRepository.SaveAdministrators(ctx, group.GetID(), []values.TraPMemberID{user.GetID()})
		if err != nil {
			return fmt.Errorf("failed to save administrators: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return &service.GroupDetail{
		Group:        group,
		Administers:  []*service.UserInfo{user},
		MainResource: mainResourceInfo,
	}, nil
}

func (g *Group) EditGroup(
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
) (*service.GroupDetail, error) {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	users, err := g.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	userMap := make(map[values.TraPMemberID]*service.UserInfo)
	for _, user := range users {
		userMap[user.GetID()] = user
	}

	var (
		group            *domain.Group
		administrators   []*service.UserInfo
		mainResourceInfo *service.ResourceInfo
	)
	err = g.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		groupInfo, err := g.groupRepository.GetGroup(ctx, id, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoGroup
		}
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		group = groupInfo.Group

		if !group.IsValidPermission() {
			return service.ErrInvalidPermission
		}

		administratorIDs, err := g.administratorRepository.GetAdministrators(ctx, groupInfo.GetID())
		if err != nil {
			return fmt.Errorf("failed to get administrators: %w", err)
		}

		for i, administrator := range administratorIDs {
			if administrator == user.GetID() {
				break
			}

			if i == len(administratorIDs)-1 {
				return service.ErrForbidden
			}
		}

		administrators = make([]*service.UserInfo, 0, len(administratorIDs))
		for _, administratorID := range administratorIDs {
			administrators = append(administrators, userMap[administratorID])
		}

		resourceInfo, err := g.resourceRepository.GetResource(ctx, mainResource)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoResource
		}
		if err != nil {
			return fmt.Errorf("failed to get main resource: %w", err)
		}

		creator, ok := userMap[resourceInfo.Creator]
		if !ok {
			return service.ErrNoUser
		}

		mainResourceInfo = &service.ResourceInfo{
			Resource: resourceInfo.Resource,
			File:     resourceInfo.File,
			Creator:  creator,
		}

		resourceList, err := g.resourceRepository.GetResourcesByIDs(ctx, resources, repository.LockTypeNone)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}
		if len(resourceList) != len(resources) {
			return service.ErrNoResource
		}

		if groupInfo.Group.GetName() != name {
			groupInfo.Group.SetName(name)
		}
		if groupInfo.Group.GetDescription() != description {
			groupInfo.Group.SetDescription(description)
		}
		if groupInfo.Group.GetType() != groupType {
			groupInfo.Group.SetType(groupType)
		}
		if groupInfo.Group.GetReadPermission() != readPermission {
			groupInfo.Group.SetReadPermission(readPermission)
		}
		if groupInfo.Group.GetWritePermission() != writePermission {
			groupInfo.Group.SetWritePermission(writePermission)
		}

		err = g.groupRepository.EditGroup(ctx, groupInfo.Group, mainResource)
		if err != nil && !errors.Is(err, repository.ErrNoRecordUpdated) {
			return fmt.Errorf("failed to save group: %w", err)
		}

		nowResources, err := g.resourceRepository.GetResources(ctx, &repository.ResourceSearchParams{
			Groups: []*domain.Group{groupInfo.Group},
		})
		if err != nil {
			return fmt.Errorf("failed to get resource: %w", err)
		}

		resourceMap := make(map[values.ResourceID]struct{})
		for _, resource := range nowResources {
			resourceMap[resource.Resource.GetID()] = struct{}{}
		}

		addResourceIDs := []values.ResourceID{}
		for _, resourceID := range resources {
			_, ok := resourceMap[resourceID]
			if ok {
				delete(resourceMap, resourceID)
			} else {
				addResourceIDs = append(addResourceIDs, resourceID)
			}
		}

		deleteResourceIDs := make([]values.ResourceID, 0, len(resourceMap))
		for resourceID := range resourceMap {
			deleteResourceIDs = append(deleteResourceIDs, resourceID)
		}

		err = g.groupRepository.AddResources(ctx, groupInfo.Group, addResourceIDs)
		if err != nil {
			return fmt.Errorf("failed to add resources: %w", err)
		}

		err = g.groupRepository.DeleteResources(ctx, groupInfo.Group, deleteResourceIDs)
		if err != nil {
			return fmt.Errorf("failed to delete resources: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return &service.GroupDetail{
		Group:        group,
		Administers:  administrators,
		MainResource: mainResourceInfo,
	}, nil
}

func (g *Group) DeleteGroup(ctx context.Context, session *domain.OIDCSession, id values.GroupID) error {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = g.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		groupInfo, err := g.groupRepository.GetGroup(ctx, id, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoGroup
		}
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		administratorIDs, err := g.administratorRepository.GetAdministrators(ctx, groupInfo.GetID())
		if err != nil {
			return fmt.Errorf("failed to get administrators: %w", err)
		}

		for i, administrator := range administratorIDs {
			if administrator == user.GetID() {
				break
			}

			if i == len(administratorIDs)-1 {
				return service.ErrForbidden
			}
		}

		err = g.groupRepository.DeleteGroup(ctx, groupInfo.Group)
		if err != nil {
			return fmt.Errorf("failed to delete group: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func (g *Group) AddResource(ctx context.Context, session *domain.OIDCSession, id values.GroupID, resource values.ResourceID) ([]*service.ResourceInfo, error) {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	users, err := g.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	userMap := make(map[values.TraPMemberID]*service.UserInfo)
	for _, user := range users {
		userMap[user.GetID()] = user
	}

	var resources []*service.ResourceInfo
	err = g.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		groupInfo, err := g.groupRepository.GetGroup(ctx, id, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoGroup
		}
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		resourceInfo, err := g.resourceRepository.GetResource(ctx, resource)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoResource
		}
		if err != nil {
			return fmt.Errorf("failed to get main resource: %w", err)
		}

		creator, ok := userMap[resourceInfo.Creator]
		if !ok {
			return service.ErrNoUser
		}

		if groupInfo.Group.GetWritePermission() != values.GroupWritePermissionPublic {
			administratorIDs, err := g.administratorRepository.GetAdministrators(ctx, groupInfo.GetID())
			if err != nil {
				return fmt.Errorf("failed to get administrators: %w", err)
			}

			for i, administrator := range administratorIDs {
				if administrator == user.GetID() {
					break
				}

				if i == len(administratorIDs)-1 {
					return service.ErrForbidden
				}
			}
		}

		nowResources, err := g.resourceRepository.GetResources(ctx, &repository.ResourceSearchParams{
			Groups: []*domain.Group{groupInfo.Group},
		})
		if err != nil {
			return fmt.Errorf("failed to get resource: %w", err)
		}

		for _, oldResource := range nowResources {
			if oldResource.Resource.GetID() == resource {
				return service.ErrResourceAlreadyExists
			}
		}

		err = g.groupRepository.AddResources(ctx, groupInfo.Group, []values.ResourceID{resource})
		if err != nil {
			return fmt.Errorf("failed to add resource: %w", err)
		}

		resources = make([]*service.ResourceInfo, 0, len(nowResources)+1)
		resources = append(resources, &service.ResourceInfo{
			Resource: resourceInfo.Resource,
			File:     resourceInfo.File,
			Creator:  creator,
		})
		for _, oldResource := range nowResources {
			resources = append(resources, &service.ResourceInfo{
				Resource: oldResource.Resource,
				File:     oldResource.File,
				Creator:  userMap[oldResource.Creator],
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return resources, nil
}

func (g *Group) GetGroup(ctx context.Context, session *domain.OIDCSession, groupID values.GroupID) (*service.GroupDetail, error) {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	users, err := g.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	userMap := make(map[values.TraPMemberID]*service.UserInfo)
	for _, user := range users {
		userMap[user.GetID()] = user
	}

	groupInfo, err := g.groupRepository.GetGroup(ctx, groupID, repository.LockTypeRecord)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoGroup
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	administratorIDs, err := g.administratorRepository.GetAdministrators(ctx, groupInfo.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	if groupInfo.Group.GetReadPermission() != values.GroupReadPermissionPublic {
		for i, administrator := range administratorIDs {
			if administrator == user.GetID() {
				break
			}

			if i == len(administratorIDs)-1 {
				return nil, service.ErrForbidden
			}
		}
	}

	administrators := make([]*service.UserInfo, 0, len(administratorIDs))
	for _, administrator := range administratorIDs {
		administrators = append(administrators, userMap[administrator])
	}

	mainResource := &service.ResourceInfo{
		Resource: groupInfo.MainResource.Resource,
		File:     groupInfo.MainResource.File,
		Creator:  userMap[groupInfo.MainResource.Creator],
	}

	return &service.GroupDetail{
		Group:        groupInfo.Group,
		Administers:  administrators,
		MainResource: mainResource,
	}, nil
}

func (g *Group) GetGroups(ctx context.Context, session *domain.OIDCSession, params *service.GroupSearchParams) ([]*service.GroupInfo, error) {
	user, err := g.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	users, err := g.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	userMap := make(map[values.TraPMemberID]*service.UserInfo)
	for _, user := range users {
		userMap[user.GetID()] = user
	}

	userNameMap := make(map[values.TraPMemberName]*service.UserInfo)
	for _, user := range users {
		userNameMap[user.GetName()] = user
	}

	userList := make([]*service.UserInfo, 0, len(params.Users))
	for _, userName := range params.Users {
		user, ok := userNameMap[userName]
		if !ok {
			return nil, service.ErrNoUser
		}

		userList = append(userList, user)
	}

	groups, err := g.groupRepository.GetGroups(ctx, user, &repository.GroupSearchParams{
		GroupTypes: params.GroupTypes,
		Users:      userList,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	groupList := make([]*service.GroupInfo, 0, len(groups))
	for _, group := range groups {
		groupList = append(groupList, &service.GroupInfo{
			Group: group.Group,
			MainResource: &service.ResourceInfo{
				Resource: group.MainResource.Resource,
				File:     group.MainResource.File,
				Creator:  userMap[group.MainResource.Creator],
			},
		})
	}

	return groupList, nil
}
