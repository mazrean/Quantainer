package v1

import (
	"github.com/mazrean/Quantainer/repository"
)

type Group struct {
	dbRepository            repository.DB
	resourceRepository      repository.Resource
	groupRepository         repository.Group
	administratorRepository repository.Administrator
	userUtils               *UserUtils
}

func NewGroup(
	dbRepository repository.DB,
	resourceRepository repository.Resource,
	groupRepository repository.Group,
	administratorRepository repository.Administrator,
	userUtils *UserUtils,
) *Group {
	return &Group{
		dbRepository:            dbRepository,
		resourceRepository:      resourceRepository,
		groupRepository:         groupRepository,
		administratorRepository: administratorRepository,
		userUtils:               userUtils,
	}
}
