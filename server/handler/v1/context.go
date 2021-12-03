package v1

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	sessionContextKey  = "session"
)

func getSession(c echo.Context) (*sessions.Session, error) {
	iSession := c.Get(sessionContextKey)
	if iSession == nil {
		return nil, errors.New("session is not set")
	}

	session, ok := iSession.(*sessions.Session)
	if !ok {
		return nil, errors.New("invalid session")
	}

	return session, nil
}
