package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/Quantainer/domain/values"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
	"github.com/mazrean/Quantainer/service"
)

type Resource struct {
	session         *Session
	checker         *Checker
	resourceService service.Resource
}

func NewResource(
	session *Session,
	checker *Checker,
	resourceService service.Resource,
) *Resource {
	return &Resource{
		session:         session,
		checker:         checker,
		resourceService: resourceService,
	}
}

func (r *Resource) PostResource(c echo.Context, strFileID Openapi.FileIDInPath) error {
	err := r.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	authSession, err := r.session.getAuthSession(session)
	if err != nil {
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get auth session")
	}

	uuidFileID, err := uuid.Parse(string(strFileID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	var newResource Openapi.NewResource
	err = c.Bind(&newResource)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	var resourceType values.ResourceType
	switch newResource.ResourceType {
	case Openapi.ResourceTypeImage:
		resourceType = values.ResourceTypeImage
	case Openapi.ResourceTypeOther:
		resourceType = values.ResourceTypeOther
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resource type")
	}

	resource, err := r.resourceService.CreateResource(
		c.Request().Context(),
		authSession,
		values.NewFileIDFromUUID(uuidFileID),
		values.NewResourceName(newResource.Name),
		resourceType,
		values.NewResourceComment(newResource.Comment),
	)
	if errors.Is(err, service.ErrNoFile) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "you are not the file owner")
	}
	if errors.Is(err, service.ErrInvalidResourceType) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resource type")
	}
	if err != nil {
		log.Printf("error: failed to create resource: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create resource")
	}

	return c.JSON(http.StatusCreated, &Openapi.Resource{
		Id:          uuid.UUID(resource.Resource.GetID()).String(),
		Creator:     string(resource.Creator.GetName()),
		FileID:      string(strFileID),
		CreatedAt:   resource.Resource.GetCreatedAt(),
		NewResource: newResource,
	})
}

func (r *Resource) GetResource(c echo.Context, resourceID Openapi.ResourceIDInPath) error {
	err := r.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	authSession, err := r.session.getAuthSession(session)
	if err != nil {
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get auth session")
	}

	uuidResourceID, err := uuid.Parse(string(resourceID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resource id")
	}

	resource, err := r.resourceService.GetResource(
		c.Request().Context(),
		authSession,
		values.NewResourceIDFromUUID(uuidResourceID),
	)
	if errors.Is(err, service.ErrNoResource) {
		return echo.NewHTTPError(http.StatusNotFound, "resource not found")
	}
	if err != nil {
		log.Printf("error: failed to get resource: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get resource")
	}

	var resourceType Openapi.ResourceType
	switch resource.Resource.GetType() {
	case values.ResourceTypeImage:
		resourceType = Openapi.ResourceTypeImage
	case values.ResourceTypeOther:
		resourceType = Openapi.ResourceTypeOther
	default:
		log.Printf("error: unknown resource type: %v\n", resource.Resource.GetType())
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid resource type")
	}

	return c.JSON(http.StatusOK, &Openapi.Resource{
		Id:        uuid.UUID(resource.Resource.GetID()).String(),
		Creator:   string(resource.Creator.GetName()),
		FileID:    uuid.UUID(resource.File.GetID()).String(),
		CreatedAt: resource.Resource.GetCreatedAt(),
		NewResource: Openapi.NewResource{
			Name:         string(resource.Resource.GetName()),
			Comment:      string(resource.Resource.GetComment()),
			ResourceType: resourceType,
		},
	})
}

func (r *Resource) GetResources(c echo.Context, params Openapi.GetResourcesParams) error {
	err := r.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	authSession, err := r.session.getAuthSession(session)
	if err != nil {
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get auth session")
	}

	var limit int
	if params.Limit != nil {
		limit = int(*params.Limit)
	} else {
		limit = -1
	}

	var offset int
	if params.Offset != nil {
		offset = int(*params.Offset)
	} else {
		offset = 0
	}

	var resourceTypes []values.ResourceType
	if params.Type != nil {
		resourceTypes = make([]values.ResourceType, 0, len(*params.Type))
		for _, resourceType := range *params.Type {
			switch resourceType {
			case Openapi.ResourceTypeImage:
				resourceTypes = append(resourceTypes, values.ResourceTypeImage)
			case Openapi.ResourceTypeOther:
				resourceTypes = append(resourceTypes, values.ResourceTypeOther)
			default:
				return echo.NewHTTPError(http.StatusBadRequest, "invalid resource type")
			}
		}
	} else {
		resourceTypes = nil
	}

	var users []values.TraPMemberName
	if params.User != nil {
		users = make([]values.TraPMemberName, 0, len(*params.User))
		for _, user := range *params.User {
			users = append(users, values.NewTrapMemberName(user))
		}
	} else {
		users = nil
	}

	resourceInfos, err := r.resourceService.GetResources(
		c.Request().Context(),
		authSession,
		&service.ResourceSearchParams{
			ResourceTypes: resourceTypes,
			Users:         users,
			Limit:         limit,
			Offset:        offset,
		},
	)
	if errors.Is(err, service.ErrNoUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user")
	}
	if err != nil {
		log.Printf("error: failed to get resources: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get resources")
	}

	resources := make([]Openapi.Resource, 0, len(resourceInfos))
	for _, resourceInfo := range resourceInfos {
		var resourceType Openapi.ResourceType
		switch resourceInfo.Resource.GetType() {
		case values.ResourceTypeImage:
			resourceType = Openapi.ResourceTypeImage
		case values.ResourceTypeOther:
			resourceType = Openapi.ResourceTypeOther
		default:
			log.Printf("error: unknown resource type: %v\n", resourceInfo.Resource.GetType())
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid resource type")
		}

		resources = append(resources, Openapi.Resource{
			Id:        uuid.UUID(resourceInfo.Resource.GetID()).String(),
			Creator:   string(resourceInfo.Creator.GetName()),
			FileID:    uuid.UUID(resourceInfo.File.GetID()).String(),
			CreatedAt: resourceInfo.Resource.GetCreatedAt(),
			NewResource: Openapi.NewResource{
				Name:         string(resourceInfo.Resource.GetName()),
				Comment:      string(resourceInfo.Resource.GetComment()),
				ResourceType: resourceType,
			},
		})
	}

	return c.JSON(http.StatusOK, resources)
}
