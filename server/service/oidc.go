package service

//go:generate mockgen -source=$GOFILE -destination=mock/${GOFILE} -package=mock

import (
	"context"

	"github.com/mazrean/Quantainer/domain"
	"github.com/mazrean/Quantainer/domain/values"
)

type OIDC interface {
	Authorize(ctx context.Context) (*domain.OIDCClient, *domain.OIDCAuthState, error)
	Callback(ctx context.Context, authState *domain.OIDCAuthState, code values.OIDCAuthorizationCode) (*domain.OIDCSession, error)
	Logout(ctx context.Context, session *domain.OIDCSession) error
	TraPAuth(ctx context.Context, session *domain.OIDCSession) error
}
