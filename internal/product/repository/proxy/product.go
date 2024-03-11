package proxy

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/pkg/cache"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/product/repository/pgrepo"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/strjoin"
	"strconv"
)

type productRepositoryImpl struct {
	pgRepo pgrepo.ProductRepository
	lc     cache.LocalCache
	logger logger.Logger
}

func NewProductRepository(pgRepo pgrepo.ProductRepository, lc cache.LocalCache, logger logger.Logger) domain.ProductRepository {
	return &productRepositoryImpl{
		pgRepo: pgRepo,
		lc:     lc,
		logger: logger,
	}
}

func (r *productRepositoryImpl) CheckProduct(ctx context.Context, productID uint64, quantity uint64) (*valueobject.ProductStatus, error) {
	status := &valueobject.ProductStatus{}
	key := strjoin.Join(checkProductKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, status)
	if ok && err == nil {
		return status, nil
	}

	status, err = r.pgRepo.CheckProduct(ctx, productID, quantity)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, status); err != nil {
		r.logger.Error("Product: failed to set product status to local cache", err)
	}

	return status, nil
}

func (r *productRepositoryImpl) GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error) {
	prodDetail := &valueobject.ProductDetail{}
	key := strjoin.Join(getProductDetailKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, prodDetail)
	if ok && err == nil {
		return prodDetail, nil
	}

	prodDetail, err = r.pgRepo.GetProductDetail(ctx, productID)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, prodDetail); err != nil {
		r.logger.Error("Product: failed to set product detail to local cache", err)
	}

	return prodDetail, nil
}

func (r *productRepositoryImpl) GetProductInventory(ctx context.Context, productID uint64) (uint64, error) {
	var prodInventory uint64
	key := strjoin.Join(getProductInventoryKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, &prodInventory)
	if ok && err == nil {
		return prodInventory, nil
	}

	prodInventory, err = r.pgRepo.GetProductInventory(ctx, productID)
	if err != nil {
		return 0, err
	}

	if err = r.lc.Set(key, prodInventory); err != nil {
		r.logger.Error("Product: failed to set product inventory to local cache", err)
	}

	return prodInventory, nil
}

func (r *productRepositoryImpl) GetProduct(ctx context.Context, productID uint64) (*entity.Product, error) {
	product := &entity.Product{}
	key := strjoin.Join(getProductKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, product)
	if ok && err == nil {
		return product, nil
	}

	product, err = r.pgRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, product); err != nil {
		r.logger.Error("Product: failed to set product to local cache", err)
	}

	return product, nil
}

func (r *productRepositoryImpl) ListProducts(ctx context.Context, limit, offset uint64) (*[]valueobject.ProductCatalog, error) {
	return r.pgRepo.ListProducts(ctx, limit, offset)
}

func (r *productRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.pgRepo.CreateProduct(ctx, product)
}

func (r *productRepositoryImpl) UpdateProductDetail(ctx context.Context, productID uint64, product *valueobject.ProductDetail) error {
	return r.pgRepo.UpdateProductDetail(ctx, productID, product)
}

func (r *productRepositoryImpl) UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	return r.pgRepo.UpdateProductInventory(ctx, idempotencyKey, purchasedProducts)
}

func (r *productRepositoryImpl) RollbackProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	return r.pgRepo.RollbackProductInventory(ctx, idempotencyKey, purchasedProducts)
}
