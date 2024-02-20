package pg_repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/product/infrastructure/db/postgres/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
)

type ProductInventory struct {
	Inventory uint64
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) CheckProduct(ctx context.Context, productID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Product{}).Where("id = ?", productID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *productRepositoryImpl) GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error) {
	var product valueobject.ProductDetail
	if err := r.db.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductInventory(ctx context.Context, productID uint64) (uint64, error) {
	var product ProductInventory
	if err := r.db.Where("id = ?", productID).Select("inventory").First(&product).Error; err != nil {
		return 0, err
	}
	return product.Inventory, nil
}

func (r *productRepositoryImpl) GetProduct(ctx context.Context, productIDs uint64) (*entity.Product, error) {
	var product model.Product
	if err := r.db.Where("id = ?", productIDs).First(&product).Error; err != nil {
		return nil, err
	}

	return &entity.Product{
		ID:         product.ID,
		CategoryID: product.CategoryID,
		Detail: &valueobject.ProductDetail{
			Name:        product.Name,
			Description: product.Description,
			BrandName:   product.BrandName,
			Price:       product.Price,
		},
		Inventory: product.Inventory,
	}, nil
}

func (r *productRepositoryImpl) ListProducts(ctx context.Context, limit, offset uint64) (*[]valueobject.ProductCatalog, error) {
	var products []valueobject.ProductCatalog
	if err := r.db.Model(&model.Product{}).
		Select("id", "category_id", "name", "price", "inventory").
		Offset(int(offset)).Limit(int(limit)).
		Find(&products).Error; err != nil {
		return nil, err
	}

	productCatalogs := make([]valueobject.ProductCatalog, len(products))
	for i, product := range products {
		productCatalogs[i] = valueobject.ProductCatalog{
			ID:         product.ID,
			CategoryID: product.CategoryID,
			Name:       product.Name,
			Price:      product.Price,
		}
	}

	return &productCatalogs, nil
}

func (r *productRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product) error {
	if err := r.db.Create(&model.Product{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Detail.Name,
		Description: product.Detail.Description,
		BrandName:   product.Detail.BrandName,
		Inventory:   product.Inventory,
		Price:       product.Detail.Price,
	}).WithContext(ctx).Error; err != nil {
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			if pgError.Code == "23505" {
				return ErrDuplicateEntry
			}
		}
		return err
	}
	return nil
}

func (r *productRepositoryImpl) UpdateProductDetail(ctx context.Context, productID uint64, product *valueobject.ProductDetail) error {
	if err := r.db.Model(&model.Product{}).Where("id = ?", productID).Updates(model.Product{
		Name:        product.Name,
		Description: product.Description,
		BrandName:   product.BrandName,
		Price:       product.Price,
	}).WithContext(ctx).Error; err != nil {
		return err
	}
	return nil
}

// UpdateProductInventory updates the inventory of purchased products
func (r *productRepositoryImpl) UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	// Update each product's inventory in purchased products in a transaction
	// Note that the idempotency key is purchased id
	// First, we need to check the idempotency key not exists in the idempotencies table
	// We start a transaction to update the inventory of each product and insert the idempotency key with product id
	// and quantity to the idempotencies table

	// Check if the idempotency is exists
	var count int64
	if err := r.db.Model(&model.Idempotency{}).Where("id = ?", idempotencyKey).Count(&count).WithContext(ctx).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrInvalidIdempotency
	}

	// sort the products by product id to avoid deadlock
	sort.Slice(*purchasedProducts, func(i, j int) bool {
		return (*purchasedProducts)[i].ID < (*purchasedProducts)[j].ID
	})

	// Start a transaction
	// With read committed isolation level and update lock, we can avoid lost update
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// Update each product's inventory
	for _, purchasedProduct := range *purchasedProducts {
		var inventory ProductInventory
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.Product{}).
			Where("id = ?", purchasedProduct.ID).Select("inventory").First(&inventory).Error; err != nil {
			tx.Rollback()
			return err
		}

		if inventory.Inventory < purchasedProduct.Quantity {
			tx.Rollback()
			return ErrInsufficientInventory
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", purchasedProduct.ID).Update("inventory", gorm.Expr("inventory - ?", purchasedProduct.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert the idempotency key to the idempotencies table
	var idempotencies []model.Idempotency
	for _, purchasedProduct := range *purchasedProducts {
		idempotencies = append(idempotencies, model.Idempotency{
			ID:         idempotencyKey,
			ProductID:  purchasedProduct.ID,
			Quantity:   purchasedProduct.Quantity,
			Rollbacked: false,
		})
	}
	if err := tx.Create(&idempotencies).WithContext(ctx).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RollbackProductInventory rollbacks the inventory of purchased products
func (r *productRepositoryImpl) RollbackProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	// Get all idempotencies
	var idempotencies []model.Idempotency
	if err := r.db.Model(&model.Idempotency{}).Select("product_id", "quantity", "rollbacked").
		Where("id = ?", idempotencyKey).Order("product_id").
		Find(&idempotencies).WithContext(ctx).Error; err != nil {
		return err
	}

	// Start a transaction
	// With read committed isolation level and update lock, we can avoid lost update
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// Rollback each product's inventory
	for _, idempotency := range idempotencies {
		var inventory ProductInventory
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.Product{}).
			Where("id = ?", idempotency.ProductID).Select("inventory").
			First(&inventory).Error; err != nil {
			tx.Rollback()
			return err
		}

		if idempotency.Rollbacked {
			continue
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", idempotency.ProductID).
			Update("inventory", gorm.Expr("inventory + ?", idempotency.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark the idempotency as rollbacked
	if err := tx.Model(&model.Idempotency{}).Where("id = ?", idempotencyKey).Update("rollbacked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
