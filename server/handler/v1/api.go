package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
)

type API struct {
	*User
	*OAuth2
	*Session
	*File
}

func NewAPI(
	user *User,
	oAuth2 *OAuth2,
	session *Session,
	file *File,
) *API {
	return &API{
		User:    user,
		OAuth2:  oAuth2,
		Session: session,
		File:    file,
	}
}

func (a *API) Start(addr string) error {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	Openapi.RegisterHandlersWithBaseURL(e, a, "/api/v1")

	return e.Start(addr)
}
