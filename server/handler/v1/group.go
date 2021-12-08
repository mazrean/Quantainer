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

	var apiGroup Openapi.NewGroup
	err = c.Bind(&apiGroup)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
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
	if errors.Is(err, service.ErrInvalidPermission) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission")
	}
	if errors.Is(err, service.ErrNoResource) {
		return echo.NewHTTPError(http.StatusBadRequest, "no resource")
	}
	if errors.Is(err, service.ErrNoUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "no user")
	}
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

func (g *Group) GetGroups(c echo.Context, params Openapi.GetGroupsParams) error {
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

	var groupTypes []values.GroupType
	if params.Type != nil {
		groupTypes = make([]values.GroupType, 0, len(*params.Type))
		for _, groupType := range *params.Type {
			switch groupType {
			case Openapi.GroupTypeArtBook:
				groupTypes = append(groupTypes, values.GroupTypeArtBook)
			case Openapi.GroupTypeOther:
				groupTypes = append(groupTypes, values.GroupTypeOther)
			default:
				return echo.NewHTTPError(http.StatusBadRequest, "invalid group type")
			}
		}
	}

	var users []values.TraPMemberName
	if params.User != nil {
		users = make([]values.TraPMemberName, 0, len(*params.User))
		for _, user := range *params.User {
			users = append(users, values.NewTrapMemberName(user))
		}
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

	groupInfos, err := g.groupServer.GetGroups(
		c.Request().Context(),
		authSession,
		&service.GroupSearchParams{
			GroupTypes: groupTypes,
			Users:      users,
			Limit:      limit,
			Offset:     offset,
		},
	)
	if errors.Is(err, service.ErrNoUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user")
	}
	if err != nil {
		log.Printf("error: failed to get groups: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get groups")
	}

	apiGroups := make([]Openapi.GroupInfo, 0, len(groupInfos))
	for _, groupInfo := range groupInfos {
		var groupType Openapi.GroupType
		switch groupInfo.Group.GetType() {
		case values.GroupTypeArtBook:
			groupType = Openapi.GroupTypeArtBook
		case values.GroupTypeOther:
			groupType = Openapi.GroupTypeOther
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid group type")
		}

		var readPermission Openapi.ReadPermission
		switch groupInfo.Group.GetReadPermission() {
		case values.GroupReadPermissionPublic:
			readPermission = Openapi.ReadPermissionPublic
		case values.GroupReadPermissionPrivate:
			readPermission = Openapi.ReadPermissionPrivate
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid group read permission")
		}

		var writePermission Openapi.WritePermission
		switch groupInfo.Group.GetWritePermission() {
		case values.GroupWritePermissionPublic:
			writePermission = Openapi.WritePermissionPublic
		case values.GroupWritePermissionPrivate:
			writePermission = Openapi.WritePermissionPrivate
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid group write permission")
		}

		var resourceType Openapi.ResourceType
		switch groupInfo.MainResource.Resource.GetType() {
		case values.ResourceTypeImage:
			resourceType = Openapi.ResourceTypeImage
		case values.ResourceTypeOther:
			resourceType = Openapi.ResourceTypeOther
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid resource type")
		}

		apiGroups = append(apiGroups, Openapi.GroupInfo{
			GroupBase: Openapi.GroupBase{
				Name:            string(groupInfo.Group.GetName()),
				Description:     string(groupInfo.Group.GetDescription()),
				Type:            groupType,
				ReadPermission:  readPermission,
				WritePermission: writePermission,
			},
			MainResource: Openapi.Resource{
				Id:        uuid.UUID(groupInfo.MainResource.Resource.GetID()).String(),
				FileID:    uuid.UUID(groupInfo.MainResource.File.GetID()).String(),
				Creator:   string(groupInfo.MainResource.Creator.GetName()),
				CreatedAt: groupInfo.GetCreatedAt(),
				NewResource: Openapi.NewResource{
					Name:         string(groupInfo.MainResource.GetName()),
					Comment:      string(groupInfo.MainResource.GetComment()),
					ResourceType: resourceType,
				},
			},
		})
	}

	return c.JSON(http.StatusOK, apiGroups)
}

func (g *Group) GetGroup(c echo.Context, strGroupID Openapi.GroupIDInPath) error {
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

	uuidGroupID, err := uuid.Parse(string(strGroupID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group id")
	}

	groupInfo, err := g.groupServer.GetGroup(
		c.Request().Context(),
		authSession,
		values.NewGroupIDFromUUID(uuidGroupID),
	)
	if errors.Is(err, service.ErrNoGroup) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group id")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	if err != nil {
		log.Printf("error: failed to get group: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get group")
	}

	var groupType Openapi.GroupType
	switch groupInfo.Group.GetType() {
	case values.GroupTypeArtBook:
		groupType = Openapi.GroupTypeArtBook
	case values.GroupTypeOther:
		groupType = Openapi.GroupTypeOther
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid group type")
	}

	var readPermission Openapi.ReadPermission
	switch groupInfo.Group.GetReadPermission() {
	case values.GroupReadPermissionPublic:
		readPermission = Openapi.ReadPermissionPublic
	case values.GroupReadPermissionPrivate:
		readPermission = Openapi.ReadPermissionPrivate
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid group read permission")
	}

	var writePermission Openapi.WritePermission
	switch groupInfo.Group.GetWritePermission() {
	case values.GroupWritePermissionPublic:
		writePermission = Openapi.WritePermissionPublic
	case values.GroupWritePermissionPrivate:
		writePermission = Openapi.WritePermissionPrivate
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid group write permission")
	}

	var resourceType Openapi.ResourceType
	switch groupInfo.MainResource.Resource.GetType() {
	case values.ResourceTypeImage:
		resourceType = Openapi.ResourceTypeImage
	case values.ResourceTypeOther:
		resourceType = Openapi.ResourceTypeOther
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid resource type")
	}

	return c.JSON(http.StatusOK, Openapi.GroupInfo{
		GroupBase: Openapi.GroupBase{
			Name:            string(groupInfo.Group.GetName()),
			Description:     string(groupInfo.Group.GetDescription()),
			Type:            groupType,
			ReadPermission:  readPermission,
			WritePermission: writePermission,
		},
		MainResource: Openapi.Resource{
			Id:        uuid.UUID(groupInfo.MainResource.Resource.GetID()).String(),
			FileID:    uuid.UUID(groupInfo.MainResource.File.GetID()).String(),
			Creator:   string(groupInfo.MainResource.Creator.GetName()),
			CreatedAt: groupInfo.GetCreatedAt(),
			NewResource: Openapi.NewResource{
				Name:         string(groupInfo.MainResource.GetName()),
				Comment:      string(groupInfo.MainResource.GetComment()),
				ResourceType: resourceType,
			},
		},
	})
}

func (g *Group) PatchGroup(c echo.Context, strGroupID Openapi.GroupIDInPath) error {
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

	uuidGroupID, err := uuid.Parse(string(strGroupID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group id")
	}

	var apiGroup Openapi.NewGroup
	err = c.Bind(&apiGroup)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
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

	groupDetail, err := g.groupServer.EditGroup(
		c.Request().Context(),
		authSession,
		values.NewGroupIDFromUUID(uuidGroupID),
		values.NewGroupName(apiGroup.Name),
		groupType,
		values.NewGroupDescription(apiGroup.Description),
		readPermission,
		writePermission,
		values.NewResourceIDFromUUID(uuidResourceID),
		resourceIDs,
	)
	if errors.Is(err, service.ErrNoGroup) {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}
	if errors.Is(err, service.ErrInvalidPermission) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission")
	}
	if errors.Is(err, service.ErrNoResource) {
		return echo.NewHTTPError(http.StatusBadRequest, "no resource")
	}
	if errors.Is(err, service.ErrNoUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "no user")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	if err != nil {
		log.Printf("error: failed to edit group: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to edit group")
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

	return c.JSON(http.StatusOK, &Openapi.GroupDetail{
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

func (g *Group) DeleteGroup(c echo.Context, strGroupID Openapi.GroupIDInPath) error {
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

	uuidGroupID, err := uuid.Parse(string(strGroupID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group id")
	}

	err = g.groupServer.DeleteGroup(
		c.Request().Context(),
		authSession,
		values.NewGroupIDFromUUID(uuidGroupID),
	)
	if errors.Is(err, service.ErrNoGroup) {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	if err != nil {
		log.Printf("error: failed to delete group: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete group")
	}

	return c.NoContent(http.StatusOK)
}

func (g *Group) PostResourceToGroup(c echo.Context, strGroupID Openapi.GroupIDInPath, strResourceID Openapi.ResourceIDInPath) error {
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

	uuidGroupID, err := uuid.Parse(string(strGroupID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid group id")
	}

	uuidResourceID, err := uuid.Parse(string(strResourceID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid resource id")
	}

	resourceInfos, err := g.groupServer.AddResource(
		c.Request().Context(),
		authSession,
		values.NewGroupIDFromUUID(uuidGroupID),
		values.NewResourceIDFromUUID(uuidResourceID),
	)
	if errors.Is(err, service.ErrNoGroup) {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}
	if errors.Is(err, service.ErrNoResource) {
		return echo.NewHTTPError(http.StatusBadRequest, "no resource")
	}
	if errors.Is(err, service.ErrNoUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "no user")
	}
	if errors.Is(err, service.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	if err != nil {
		log.Printf("error: failed to edit group: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to edit group")
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
