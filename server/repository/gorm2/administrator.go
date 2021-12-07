package gorm2

type Administrator struct {
	db *DB
}

func NewAdministrator(db *DB) *Administrator {
	return &Administrator{
		db: db,
	}
}
