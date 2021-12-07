package gorm2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mazrean/Quantainer/domain/values"
)

type Administrator struct {
	db *DB
}

func NewAdministrator(db *DB) *Administrator {
	return &Administrator{
		db: db,
	}
}

func (a *Administrator) SaveAdministrators(ctx context.Context, groupID values.GroupID, administrators []values.TraPMemberID) error {
	db, err := a.db.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}

	administratorTables := make([]*AdministratorTable, 0, len(administrators))
	for _, administrator := range administrators {
		administratorTables = append(administratorTables, &AdministratorTable{
			GroupID: uuid.UUID(groupID),
			UserID:  uuid.UUID(administrator),
		})
	}

	err = db.Create(&administratorTables).Error
	if err != nil {
		return fmt.Errorf("failed to save administrators: %w", err)
	}

	return nil
}
