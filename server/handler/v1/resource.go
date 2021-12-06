package v1

import "github.com/mazrean/Quantainer/service"

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
