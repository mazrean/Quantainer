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
