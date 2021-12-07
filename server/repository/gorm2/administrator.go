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

func (a *Administrator) GetAdministrators(ctx context.Context, groupID values.GroupID) ([]values.TraPMemberID, error) {
	db, err := a.db.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %w", err)
	}

	var administratorTables []AdministratorTable
	err = db.
		Where("group_id = ?", uuid.UUID(groupID)).
		Find(&administratorTables).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators: %w", err)
	}

	administrators := make([]values.TraPMemberID, 0, len(administratorTables))
	for _, administratorTable := range administratorTables {
		administrators = append(administrators, values.NewTrapMemberID(administratorTable.UserID))
	}

	return administrators, nil
}
