package v1

import "github.com/mazrean/Quantainer/service"

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
