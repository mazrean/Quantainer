package values

import "github.com/google/uuid"

type (
	ResourceID      uuid.UUID
	ResourceName    string
	ResourceType    int8
	ResourceComment string
)

func NewResourceID() ResourceID {
	return ResourceID(uuid.New())
}

func NewResourceIDFromUUID(u uuid.UUID) ResourceID {
	return ResourceID(u)
}

func NewResourceName(name string) ResourceName {
	return ResourceName(name)
}

const (
	ResourceTypeImage = iota + 1
	ResourceTypeOther
)

func NewResourceComment(comment string) ResourceComment {
	return ResourceComment(comment)
}
