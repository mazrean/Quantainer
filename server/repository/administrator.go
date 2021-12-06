package repository

import (
	"context"

	"github.com/mazrean/Quantainer/domain/values"
)

type Administrator interface {
	SaveAdministrators(ctx context.Context, groupID values.GroupID, admin []values.TraPMemberID) error
	GetAdministrators(ctx context.Context, groupID values.GroupID) ([]values.TraPMemberID, error)
}
