package postgres

import (
	"github.com/scul0405/saga-orchestration/services/product/config"
	"github.com/scul0405/saga-orchestration/services/product/internal/infrastructure/db/postgres/model"
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
		if err := m.db.Migrator().DropTable(&model.Category{}); err != nil {
			return err
		}

		if err := m.db.Migrator().DropTable(&model.Product{}); err != nil {
			return err
		}

		if err := m.db.Migrator().DropTable(&model.Idempotency{}); err != nil {
			return err
		}
	}

	return m.db.AutoMigrate(&model.Category{}, &model.Product{}, &model.Idempotency{})
}
