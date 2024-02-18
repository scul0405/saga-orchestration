package postgres

import (
	"github.com/scul0405/saga-orchestration/services/order/config"
	"github.com/scul0405/saga-orchestration/services/order/internal/infrastructure/db/postgres/model"
	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(migration config.Migration) error {
	if !migration.Enable {
		return nil
	}

	if migration.Recreate {
		if err := m.db.Migrator().DropTable(&model.Order{}); err != nil {
			return err
		}
	}

	return m.db.AutoMigrate(&model.Order{})
}
