package v1

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
	"github.com/mazrean/Quantainer/service"
)

type User struct {
	session     *Session
	checker     *Checker
	userService service.User
}

func NewUser(session *Session, checker *Checker, userService service.User) *User {
	return &User{
		session:     session,
		checker:     checker,
		userService: userService,
	}
}

func (u *User) GetMe(c echo.Context) error {
	err := u.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	authSession, err := u.session.getAuthSession(session)
	if err != nil {
		// middlewareでログイン済みなことは確認しているので、ここではエラーになりえないはず
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	userInfo, err := u.userService.GetMe(c.Request().Context(), authSession)
	if err != nil {
		log.Printf("error: failed to get user info: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, Openapi.User{
		Id:   uuid.UUID(userInfo.GetID()).String(),
		Name: string(userInfo.GetName()),
	})
}

func (u *User) GetUsers(c echo.Context) error {
	err := u.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	authSession, err := u.session.getAuthSession(session)
	if err != nil {
		// middlewareでログイン済みなことは確認しているので、ここではエラーになりえないはず
		log.Printf("error: failed to get auth session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	userInfos, err := u.userService.GetAllActiveUser(c.Request().Context(), authSession)
	if err != nil {
		log.Printf("error: failed to get user info: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	users := make([]Openapi.User, 0, len(userInfos))
	for _, userInfo := range userInfos {
		users = append(users, Openapi.User{
			Id:   uuid.UUID(userInfo.GetID()).String(),
			Name: string(userInfo.GetName()),
		})
	}

	return c.JSON(http.StatusOK, users)
}
