package domain

import (
	"time"

	"github.com/mazrean/Quantainer/domain/values"
)

type Group struct {
	id              values.GroupID
	name            values.GroupName
	groupType       values.GroupType
	description     values.GroupDescription
	readPermission  values.GroupReadPermission
	writePermission values.GroupWritePermission
	createdAt       time.Time
}

func NewGroup(
	id values.GroupID,
	name values.GroupName,
	groupType values.GroupType,
	description values.GroupDescription,
	readPermission values.GroupReadPermission,
	writePermission values.GroupWritePermission,
	createdAt time.Time,
) *Group {
	return &Group{
		id:              id,
		name:            name,
		groupType:       groupType,
		description:     description,
		readPermission:  readPermission,
		writePermission: writePermission,
		createdAt:       createdAt,
	}
}

func (g *Group) GetID() values.GroupID {
	return g.id
}

func (g *Group) GetName() values.GroupName {
	return g.name
}

func (g *Group) GetType() values.GroupType {
	return g.groupType
}

func (g *Group) GetDescription() values.GroupDescription {
	return g.description
}

func (g *Group) GetReadPermission() values.GroupReadPermission {
	return g.readPermission
}

func (g *Group) GetWritePermission() values.GroupWritePermission {
	return g.writePermission
}

func (g *Group) IsValidPermission() bool {
	if g.readPermission == values.GroupReadPermissionPrivate && g.writePermission == values.GroupWritePermissionPublic {
		return false
	}

	return true
}

func (g *Group) GetCreatedAt() time.Time {
	return g.createdAt
}
