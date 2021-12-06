package v1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/repository"
	"github.com/mazrean/Quantainer/service"
	"github.com/mazrean/Quantainer/storage"
)

type Resource struct {
	dbRepository       repository.DB
	fileRepository     repository.File
	resourceRepository repository.Resource
	fileStorage        storage.File
	userUtils          *UserUtils
}

func NewResource(
	dbRepository repository.DB,
	fileRepository repository.File,
	resourceRepository repository.Resource,
	fileStorage storage.File,
	userUtils *UserUtils,
) *Resource {
	return &Resource{
		dbRepository:       dbRepository,
		fileRepository:     fileRepository,
		resourceRepository: resourceRepository,
		fileStorage:        fileStorage,
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

	var resource *domain.Resource
	err = r.dbRepository.Transaction(ctx, nil, func(ctx context.Context) error {
		fileInfo, err := r.fileRepository.GetFile(ctx, fileID, repository.LockTypeRecord)
		if errors.Is(err, repository.ErrRecordNotFound) {
			return service.ErrNoFile
		}
		if err != nil {
			return fmt.Errorf("failed to get file: %w", err)
		}

		if fileInfo.Creator != user.GetID() {
			return service.ErrForbidden
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

	resourceInfos, err := r.resourceRepository.GetResources(ctx, &repository.ResourceSearchParams{
		ResourceTypes: params.ResourceTypes,
		Users:         userList,
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
			Creator:  user,
		})
	}

	return resources, nil
}

func (r *Resource) DownloadResourceFile(ctx context.Context, resourceID values.ResourceID, writer io.Writer) (*domain.File, error) {
	file, err := r.fileRepository.GetFileByResourceID(ctx, resourceID)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, service.ErrNoResource
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	err = r.fileStorage.GetFile(ctx, file, writer)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}
