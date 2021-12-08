package v1

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/Quantainer/domain/values"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
	"github.com/mazrean/Quantainer/service"
)

type Group struct {
	session     *Session
	checker     *Checker
	groupServer service.Group
}

func NewGroup(
	session *Session,
	checker *Checker,
	groupServer service.Group,
) *Group {
	return &Group{
		session:     session,
		checker:     checker,
		groupServer: groupServer,
	}
}

func (g *Group) PostGroup(c echo.Context) error {
	var apiGroup Openapi.NewGroup
	err := c.Bind(&apiGroup)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	authSession, err := g.session.getAuthSession(session)
	if err != nil {
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get auth session")
	}

	var groupType values.GroupType
	switch apiGroup.Type {
	case Openapi.GroupTypeArtBook:
		groupType = values.GroupTypeArtBook
	case Openapi.GroupTypeOther:
		groupType = values.GroupTypeOther
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group type")
	}

	var readPermission values.GroupReadPermission
	switch apiGroup.ReadPermission {
	case Openapi.ReadPermissionPublic:
		readPermission = values.GroupReadPermissionPublic
	case Openapi.ReadPermissionPrivate:
		readPermission = values.GroupReadPermissionPrivate
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid read permission")
	}

	var writePermission values.GroupWritePermission
	switch apiGroup.WritePermission {
	case Openapi.WritePermissionPublic:
		writePermission = values.GroupWritePermissionPublic
	case Openapi.WritePermissionPrivate:
		writePermission = values.GroupWritePermissionPrivate
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid write permission")
	}

	uuidResourceID, err := uuid.Parse(apiGroup.MainResourceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resource id")
	}

	resourceIDs := make([]values.ResourceID, 0, len(apiGroup.ResourceIDs))
	for _, resourceID := range apiGroup.ResourceIDs {
		uuidResourceID, err := uuid.Parse(resourceID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid resource id")
		}
		resourceIDs = append(resourceIDs, values.ResourceID(uuidResourceID))
	}

	groupDetail, err := g.groupServer.CreateGroup(
		c.Request().Context(),
		authSession,
		values.NewGroupName(apiGroup.Name),
		groupType,
		values.NewGroupDescription(apiGroup.Description),
		readPermission,
		writePermission,
		values.NewResourceIDFromUUID(uuidResourceID),
		resourceIDs,
	)
	if err != nil {
		log.Printf("error: failed to create group: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create group")
	}

	var resourceType Openapi.ResourceType
	switch groupDetail.MainResource.Resource.GetType() {
	case values.ResourceTypeImage:
		resourceType = Openapi.ResourceTypeImage
	case values.ResourceTypeOther:
		resourceType = Openapi.ResourceTypeOther
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid resource type")
	}

	administrators := make([]string, 0, len(groupDetail.Administers))
	for _, administrator := range groupDetail.Administers {
		administrators = append(administrators, string(administrator.GetName()))
	}

	return c.JSON(http.StatusCreated, &Openapi.GroupDetail{
		GroupBase:      apiGroup.GroupBase,
		Administrators: administrators,
		MainResource: Openapi.Resource{
			Id:        uuid.UUID(groupDetail.MainResource.Resource.GetID()).String(),
			FileID:    uuid.UUID(groupDetail.MainResource.File.GetID()).String(),
			Creator:   string(groupDetail.MainResource.Creator.GetName()),
			CreatedAt: groupDetail.GetCreatedAt(),
			NewResource: Openapi.NewResource{
				Name:         string(groupDetail.MainResource.GetName()),
				Comment:      string(groupDetail.MainResource.GetComment()),
				ResourceType: resourceType,
			},
		},
	})
}
