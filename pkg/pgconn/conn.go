package pgconn

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	maxOpenConns    = 60
	maxIdleConns    = 30
	connMaxIdleTime = 20
	connMaxLifetime = 120
)

func NewPsqlDB(dataSourceName string) (*gorm.DB, error) {
	gormDb, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	})
	if err != nil {
		return nil, err
	}

	db, err := gormDb.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return gormDb, nil
}
