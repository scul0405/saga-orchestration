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
	"gorm.io/gorm"
	"strconv"
)

type productRepositoryImpl struct {
	pgRepo pgrepo.ProductRepository
	lc     cache.LocalCache
	rc     cache.RedisCache
	logger logger.Logger
}

func NewProductRepository(
	pgRepo pgrepo.ProductRepository,
	lc cache.LocalCache,
	rc cache.RedisCache,
	logger logger.Logger) (domain.ProductRepository, error) {

	exist, err := rc.CFExist(context.Background(), cuckooFilter, dummnyItem)
	if err != nil {
		logger.Error("Product: failed to check cuckoo filter existence", err)
		return nil, err
	}

	if !exist {
		if err = rc.CFReserve(context.Background(), cuckooFilter, 1000, 4, 100); err != nil {
			logger.Error("Product: failed to reserve cuckoo filter", err)
			return nil, err
		}
		if err = rc.CFAdd(context.Background(), cuckooFilter, dummnyItem); err != nil {
			logger.Error("Product: failed to add dummy item to cuckoo filter", err)
			return nil, err
		}

		logger.Info("Product: cuckoo filter created")
	} else {
		logger.Info("Product: cuckoo filter already exist")
	}

	return &productRepositoryImpl{
		pgRepo: pgRepo,
		lc:     lc,
		rc:     rc,
		logger: logger,
	}, nil
}

func (r *productRepositoryImpl) CheckProduct(ctx context.Context, productID uint64, quantity uint64) (*valueobject.ProductStatus, error) {
	status := &valueobject.ProductStatus{}
	key := strjoin.Join(checkProductKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, status)
	if ok && err == nil {
		return status, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, productID)
	r.logger.Error(err)
	if !exist && err == nil {
		return &valueobject.ProductStatus{
			ID:     productID,
			Price:  0,
			Status: false,
		}, nil
	}

	ok, err = r.rc.Get(ctx, key, status)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, status))
		return status, nil
	}

	// lock to prevent lost update
	mu := r.rc.GetMutex(strjoin.Join(mutexKey, key))
	if err = mu.Lock(); err != nil {
		return nil, err
	}
	defer mu.Unlock()

	// Get again to prevent new update
	ok, err = r.rc.Get(ctx, key, status)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, status))
		return status, nil
	}

	status, err = r.pgRepo.CheckProduct(ctx, productID, quantity)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, status); err != nil {
		r.logger.Error("Product: failed to set product status to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, status))
	return status, nil
}

func (r *productRepositoryImpl) GetProductDetail(ctx context.Context, productID uint64) (*valueobject.ProductDetail, error) {
	prodDetail := &valueobject.ProductDetail{}
	key := strjoin.Join(getProductDetailKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, prodDetail)
	if ok && err == nil {
		return prodDetail, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, productID)
	r.logger.Error(err)
	if !exist && err == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ok, err = r.rc.Get(ctx, key, prodDetail)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, prodDetail))
		return prodDetail, nil
	}

	// lock to prevent lost update
	mu := r.rc.GetMutex(strjoin.Join(mutexKey, key))
	if err = mu.Lock(); err != nil {
		return nil, err
	}
	defer mu.Unlock()

	// Get again to prevent new update
	ok, err = r.rc.Get(ctx, key, prodDetail)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, prodDetail))
		return prodDetail, nil
	}

	prodDetail, err = r.pgRepo.GetProductDetail(ctx, productID)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, prodDetail); err != nil {
		r.logger.Error("Product: failed to set product detail to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, prodDetail))
	return prodDetail, nil
}

func (r *productRepositoryImpl) GetProductInventory(ctx context.Context, productID uint64) (uint64, error) {
	var prodInventory uint64
	key := strjoin.Join(getProductInventoryKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, &prodInventory)
	if ok && err == nil {
		return prodInventory, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, productID)
	r.logger.Error(err)
	if !exist && err == nil {
		return 0, gorm.ErrRecordNotFound
	}

	ok, err = r.rc.Get(ctx, key, prodInventory)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, prodInventory))
		return prodInventory, nil
	}

	// lock to prevent lost update
	mu := r.rc.GetMutex(strjoin.Join(mutexKey, key))
	if err = mu.Lock(); err != nil {
		return 0, err
	}
	defer mu.Unlock()

	// Get again to prevent new update
	ok, err = r.rc.Get(ctx, key, prodInventory)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, prodInventory))
		return prodInventory, nil
	}

	prodInventory, err = r.pgRepo.GetProductInventory(ctx, productID)
	if err != nil {
		return 0, err
	}

	if err = r.lc.Set(key, prodInventory); err != nil {
		r.logger.Error("Product: failed to set product inventory to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, prodInventory))
	return prodInventory, nil
}

func (r *productRepositoryImpl) GetProduct(ctx context.Context, productID uint64) (*entity.Product, error) {
	product := &entity.Product{}
	key := strjoin.Join(getProductKey, strconv.FormatUint(productID, 10))
	ok, err := r.lc.Get(key, product)
	if ok && err == nil {
		return product, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, productID)
	r.logger.Error(err)
	if !exist && err == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ok, err = r.rc.Get(ctx, key, product)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, product))
		return product, nil
	}

	// lock to prevent lost update
	mu := r.rc.GetMutex(strjoin.Join(mutexKey, key))
	if err = mu.Lock(); err != nil {
		return nil, err
	}
	defer mu.Unlock()

	// Get again to prevent new update
	ok, err = r.rc.Get(ctx, key, product)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, product))
		return product, nil
	}

	product, err = r.pgRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, product); err != nil {
		r.logger.Error("Product: failed to set product to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, product))
	return product, nil
}

func (r *productRepositoryImpl) ListProducts(ctx context.Context, limit, offset uint64) (*[]valueobject.ProductCatalog, error) {
	return r.pgRepo.ListProducts(ctx, limit, offset)
}

func (r *productRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product) error {
	err := r.pgRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	r.logger.Error(r.rc.CFAdd(ctx, cuckooFilter, product.ID))
	return nil
}

func (r *productRepositoryImpl) UpdateProductDetail(ctx context.Context, productID uint64, product *valueobject.ProductDetail) error {
	err := r.pgRepo.UpdateProductDetail(ctx, productID, product)
	if err != nil {
		return err
	}

	key := strjoin.Join(getProductDetailKey, strconv.FormatUint(productID, 10))
	r.logger.Error(r.lc.Set(key, product))
	r.logger.Error(r.rc.Set(ctx, key, product))

	return nil
}

func (r *productRepositoryImpl) UpdateProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	err := r.pgRepo.UpdateProductInventory(ctx, idempotencyKey, purchasedProducts)
	if err != nil {
		return err
	}

	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
	for i, purchasedProduct := range *purchasedProducts {
		payloads[i] = cache.RedisIncrbyXPayload{
			Key:   strjoin.Join(getProductInventoryKey, strconv.FormatUint(purchasedProduct.ID, 10)),
			Value: int64(-purchasedProduct.Quantity),
		}
	}
	if len(payloads) > 0 {
		r.logger.Error(r.rc.ExecIncrbyXPipeline(ctx, &payloads))
	}

	return nil
}

func (r *productRepositoryImpl) RollbackProductInventory(ctx context.Context, idempotencyKey uint64, purchasedProducts *[]valueobject.PurchasedProduct) error {
	err := r.pgRepo.RollbackProductInventory(ctx, idempotencyKey, purchasedProducts)
	if err != nil {
		return err
	}

	payloads := make([]cache.RedisIncrbyXPayload, len(*purchasedProducts))
	for i, purchasedProduct := range *purchasedProducts {
		payloads[i] = cache.RedisIncrbyXPayload{
			Key:   strjoin.Join(getProductInventoryKey, strconv.FormatUint(purchasedProduct.ID, 10)),
			Value: int64(purchasedProduct.Quantity),
		}
	}
	if len(payloads) > 0 {
		r.logger.Error(r.rc.ExecIncrbyXPipeline(ctx, &payloads))
	}

	return nil
}
