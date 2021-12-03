package auth

//go:generate mockgen -source=$GOFILE -destination=mock/${GOFILE} -package=mock

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/service"
)

type User interface {
	GetMe(ctx context.Context, session *domain.OIDCSession) (*service.UserInfo, error)
	GetAllActiveUsers(ctx context.Context, session *domain.OIDCSession) ([]*service.UserInfo, error)
}
