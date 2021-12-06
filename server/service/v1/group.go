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
