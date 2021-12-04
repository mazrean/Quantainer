package values

import "github.com/google/uuid"

type (
	GroupID              uuid.UUID
	GroupName            string
	GroupDescription     string
	GroupType            int8
	GroupReadPermission  int8
	GroupWritePermission int8
)

func NewGroupID() GroupID {
	return GroupID(uuid.New())
}

func NewGroupIDFromUUID(u uuid.UUID) GroupID {
	return GroupID(u)
}

func NewGroupName(name string) GroupName {
	return GroupName(name)
}

func NewGroupDescription(desc string) GroupDescription {
	return GroupDescription(desc)
}

const (
	GroupTypeArtBook GroupType = iota + 1
	GroupTypeOther
)

const (
	GroupReadPermissionPublic GroupReadPermission = iota + 1
	GroupReadPermissionPrivate
)

const (
	GroupWritePermissionPublic GroupWritePermission = iota + 1
	GroupWritePermissionPrivate
)
