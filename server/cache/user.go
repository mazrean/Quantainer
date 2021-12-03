package cache

//go:generate mockgen -source=$GOFILE -destination=mock/${GOFILE} -package=mock

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
	"github.com/mazrean/Quantainer/service"
)

type User interface {
	GetMe(ctx context.Context, accessToken values.OIDCAccessToken) (*service.UserInfo, error)
	SetMe(ctx context.Context, session *domain.OIDCSession, user *service.UserInfo) error
	GetAllActiveUsers(ctx context.Context) ([]*service.UserInfo, error)
	SetAllActiveUsers(ctx context.Context, users []*service.UserInfo) error
}
