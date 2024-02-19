package postgres

import (
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	model2 "github.com/scul0405/saga-orchestration/services/product/infrastructure/db/postgres/model"
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
		if err := m.db.Migrator().DropTable(&model2.Category{}); err != nil {
			return err
		}

		if err := m.db.Migrator().DropTable(&model2.Product{}); err != nil {
			return err
		}

		if err := m.db.Migrator().DropTable(&model2.Idempotency{}); err != nil {
			return err
		}
	}

	return m.db.AutoMigrate(&model2.Category{}, &model2.Product{}, &model2.Idempotency{})
}
