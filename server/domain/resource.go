package domain

import (
	"time"

	"github.com/mazrean/Quantainer/domain/values"
)

type Resource struct {
	id           values.ResourceID
	name         values.ResourceName
	resourceType values.ResourceType
	comment      values.ResourceComment
	createdAt    time.Time
}

func NewResource(
	id values.ResourceID,
	name values.ResourceName,
	resourceType values.ResourceType,
	comment values.ResourceComment,
	createdAt time.Time,
) *Resource {
	return &Resource{
		id:           id,
		name:         name,
		resourceType: resourceType,
		comment:      comment,
		createdAt:    createdAt,
	}
}

func (r *Resource) GetID() values.ResourceID {
	return r.id
}

func (r *Resource) GetName() values.ResourceName {
	return r.name
}

func (r *Resource) GetType() values.ResourceType {
	return r.resourceType
}

func (r *Resource) GetComment() values.ResourceComment {
	return r.comment
}

func (r *Resource) GetCreatedAt() time.Time {
	return r.createdAt
}
