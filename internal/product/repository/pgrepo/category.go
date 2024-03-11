package pgrepo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/product/infrastructure/db/postgres/model"
	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepositoryImpl(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) CheckCategory(ctx context.Context, categoryID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Category{}).Where("id = ?", categoryID).Count(&count).WithContext(ctx).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepositoryImpl) GetCategory(ctx context.Context, categoryID uint64) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.Where("id = ?", categoryID).First(&category).WithContext(ctx).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) CreateCategory(ctx context.Context, category *entity.Category) error {
	if err := r.db.Create(category).WithContext(ctx).Error; err != nil {
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			if pgError.Code == "23505" {
				return ErrDuplicateEntry
			}
		}
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) UpdateCategory(ctx context.Context, categoryID uint64, category *entity.Category) error {
	return r.db.Model(&model.Category{}).Where("id = ?", categoryID).Updates(&model.Category{
		Name:        category.Name,
		Description: category.Description,
	}).WithContext(ctx).Error
}
