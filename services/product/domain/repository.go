package domain

import (
	"context"
	entity2 "github.com/scul0405/saga-orchestration/services/product/domain/entity"
	"github.com/scul0405/saga-orchestration/services/product/domain/valueobject"
)

// ProductRepository is an interface for product repository
type ProductRepository interface {
	CheckProduct(ctx context.Context, productID uint64) (bool, error)
	GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error)
	GetProductInventory(ctx context.Context, productID uint64) (uint64, error)
	GetProduct(ctx context.Context, productIDs uint64) (*entity2.Product, error)
	ListProducts(ctx context.Context, limit, offset uint64) (*[]valueobject.ProductCatalog, error)
	CreateProduct(ctx context.Context, product *entity2.Product) error
	UpdateProductDetail(ctx context.Context, productID uint64, product *valueobject.ProductDetail) error
	UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error
	RollbackProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error
}

// CategoryRepository is an interface for category repository
type CategoryRepository interface {
	CheckCategory(ctx context.Context, categoryID uint64) (bool, error)
	GetCategory(ctx context.Context, categoryID uint64) (*entity2.Category, error)
	CreateCategory(ctx context.Context, category *entity2.Category) error
	UpdateCategory(ctx context.Context, categoryID uint64, category *entity2.Category) error
}
