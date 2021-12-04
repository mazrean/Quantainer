package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
	"github.com/mazrean/Quantainer/service"
)

type Checker struct {
	session     *Session
	oidcService service.OIDC
}

func NewChecker(
	session *Session,
	oidcService service.OIDC,
) *Checker {
	return &Checker{
		session:     session,
		oidcService: oidcService,
	}
}

func (m *Checker) check(c echo.Context) error {
	checkers := []struct {
		key     string
		handler echo.HandlerFunc
	}{
		{
			key:     Openapi.TraPMemberAuthScopes,
			handler: m.TrapMemberAuthChecker,
		},
	}

	for _, checker := range checkers {
		if c.Get(checker.key) != nil {
			if err := checker.handler(c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Checker) TrapMemberAuthChecker(c echo.Context) error {
	ok, message, err := m.checkTrapMemberAuth(c)
	if err != nil {
		log.Printf("error: failed to check trap member auth: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, message)
	}

	return nil
}

func (m *Checker) checkTrapMemberAuth(c echo.Context) (bool, string, error) {
	session, err := m.session.getSession(c)
	if err != nil {
		return false, "", fmt.Errorf("failed to get session: %w", err)
	}

	authSession, err := m.session.getAuthSession(session)
	if errors.Is(err, ErrNoValue) {
		return false, "no access token", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("failed to get auth session: %w", err)
	}

	err = m.oidcService.TraPAuth(c.Request().Context(), authSession)
	if errors.Is(err, service.ErrOIDCSessionExpired) {
		return false, "access token is expired", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("failed to check traP auth: %w", err)
	}

	return true, "", nil
}
