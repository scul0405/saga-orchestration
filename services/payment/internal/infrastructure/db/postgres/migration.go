package postgres

import (
	"github.com/scul0405/saga-orchestration/services/payment/config"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/db/postgres/model"
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
		if err := m.db.Migrator().DropTable(&model.Payment{}); err != nil {
			return err
		}
	}

	return m.db.AutoMigrate(&model.Payment{})
}
