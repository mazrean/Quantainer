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

type Resource struct {
	dbRepository       repository.DB
	fileRepository     repository.File
	resourceRepository repository.Resource
	groupRepository    repository.Group
	userUtils          *UserUtils
}

func NewResource(
	dbRepository repository.DB,
	fileRepository repository.File,
	resourceRepository repository.Resource,
	groupRepository repository.Group,
	userUtils *UserUtils,
) *Resource {
	return &Resource{
		dbRepository:       dbRepository,
		fileRepository:     fileRepository,
		resourceRepository: resourceRepository,
		groupRepository:    groupRepository,
		userUtils:          userUtils,
	}
}

func (r *Resource) CreateResource(
	ctx context.Context,
	session *domain.OIDCSession,
	fileID values.FileID,
	name values.ResourceName,
	resourceType values.ResourceType,
	comment values.ResourceComment,
) (*service.ResourceInfo, error) {
	user, err := r.userUtils.getMe(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var fileInfo *repository.FileWithCreator

	var resource *domain.Resource
	err = r.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		var err error
		fileInfo, err = r.fileRepository.GetFile(ctx, fileID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoFile
		}
		if err != nil {
			return fmt.Errorf("failed to get file: %w", err)
		}

		if fileInfo.Creator != user.GetID() {
			return service.ErrForbidden
		}

		if !fileInfo.File.GetType().IsValidResourceType(resourceType) {
			return service.ErrInvalidResourceType
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

	return &service.ResourceInfo{
		Resource: resource,
		File:     fileInfo.File,
		Creator:  user,
	}, nil
}

func (r *Resource) CreateBotResource(
	ctx context.Context,
	user *service.UserInfo,
	fileID values.FileID,
	name values.ResourceName,
	resourceType values.ResourceType,
	comment values.ResourceComment,
	createdAt time.Time,
) (*service.ResourceInfo, error) {
	var fileInfo *repository.FileWithCreator

	var resource *domain.Resource
	err := r.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		var err error
		fileInfo, err = r.fileRepository.GetFile(ctx, fileID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoFile
		}
		if err != nil {
			return fmt.Errorf("failed to get file: %w", err)
		}

		if fileInfo.Creator != user.GetID() {
			return service.ErrForbidden
		}

		if !fileInfo.File.GetType().IsValidResourceType(resourceType) {
			return service.ErrInvalidResourceType
		}

		resource = domain.NewResource(
			values.NewResourceID(),
			name,
			resourceType,
			comment,
			createdAt,
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

	return &service.ResourceInfo{
		Resource: resource,
		File:     fileInfo.File,
		Creator:  user,
	}, nil
}

func (r *Resource) GetResource(ctx context.Context, session *domain.OIDCSession, resourceID values.ResourceID) (*service.ResourceInfo, error) {
	resourceInfo, err := r.resourceRepository.GetResource(ctx, resourceID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoResource
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	users, err := r.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var creator *service.UserInfo
	for i, user := range users {
		if user.GetID() == resourceInfo.Creator {
			creator = user
			break
		}

		if i == len(users)-1 {
			return nil, service.ErrNoUser
		}
	}

	return &service.ResourceInfo{
		Resource: resourceInfo.Resource,
		File:     resourceInfo.File,
		Creator:  creator,
	}, nil
}

func (r *Resource) GetResources(ctx context.Context, session *domain.OIDCSession, params *service.ResourceSearchParams) ([]*service.ResourceInfo, error) {
	users, err := r.userUtils.getAllActiveUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
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

	var groups []*domain.Group
	if params.Group != nil {
		groupInfos, err := r.groupRepository.GetGroup(ctx, *params.Group, repository.LockTypeNone)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, service.ErrNoGroup
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get groups: %w", err)
		}

		groups = []*domain.Group{groupInfos.Group}
	}

	resourceInfos, err := r.resourceRepository.GetResources(ctx, &repository.ResourceSearchParams{
		ResourceTypes: params.ResourceTypes,
		Users:         userList,
		Groups:        groups,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	userMap := make(map[values.TraPMemberID]*service.UserInfo)
	for _, user := range users {
		userMap[user.GetID()] = user
	}

	resources := make([]*service.ResourceInfo, 0, len(resourceInfos))
	for _, resourceInfo := range resourceInfos {
		user, ok := userMap[resourceInfo.Creator]
		if !ok {
			return nil, service.ErrNoUser
		}

		resources = append(resources, &service.ResourceInfo{
			Resource: resourceInfo.Resource,
			File:     resourceInfo.File,
			Creator:  user,
		})
	}

	return resources, nil
}
