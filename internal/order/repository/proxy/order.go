package proxy

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/order/domain"
	"github.com/scul0405/saga-orchestration/internal/order/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/order/repository/pgrepo"
	"github.com/scul0405/saga-orchestration/internal/pkg/cache"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/strjoin"
	"gorm.io/gorm"
	"strconv"
)

type orderRepositoryImpl struct {
	pgRepo pgrepo.OrderRepository
	lc     cache.LocalCache
	rc     cache.RedisCache
	logger logger.Logger
}

func NewOrderRepository(
	pgRepo pgrepo.OrderRepository,
	lc cache.LocalCache,
	rc cache.RedisCache,
	logger logger.Logger) (domain.OrderRepository, error) {

	exist, err := rc.CFExist(context.Background(), cuckooFilter, dummnyItem)
	if err != nil {
		logger.Error("Order: failed to check cuckoo filter existence", err)
		return nil, err
	}

	if !exist {
		if err = rc.CFReserve(context.Background(), cuckooFilter, 1000, 4, 100); err != nil {
			logger.Error("Order: failed to reserve cuckoo filter", err)
			return nil, err
		}
		if err = rc.CFAdd(context.Background(), cuckooFilter, dummnyItem); err != nil {
			logger.Error("Order: failed to add dummy item to cuckoo filter", err)
			return nil, err
		}

		logger.Info("Order: cuckoo filter created")
	} else {
		logger.Info("Order: cuckoo filter already exist")
	}

	return &orderRepositoryImpl{
		pgRepo: pgRepo,
		lc:     lc,
		rc:     rc,
		logger: logger,
	}, nil
}

func (r *orderRepositoryImpl) GetOrder(ctx context.Context, id uint64) (*entity.Order, error) {
	order := &entity.Order{}
	key := strjoin.Join(getOrderKey, strconv.FormatUint(id, 10))
	ok, err := r.lc.Get(key, order)
	if ok && err == nil {
		return order, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, id)
	r.logger.Error(err)
	if !exist && err == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ok, err = r.rc.Get(ctx, key, order)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, order))
		return order, nil
	}

	order, err = r.pgRepo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, order); err != nil {
		r.logger.Error("Order: failed to set order to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, order))
	return order, nil
}

func (r *orderRepositoryImpl) CreateOrder(ctx context.Context, order *entity.Order) error {
	err := r.pgRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	r.logger.Error(r.rc.CFAdd(ctx, cuckooFilter, order.ID))
	return nil
}

func (r *orderRepositoryImpl) DeleteOrder(ctx context.Context, id uint64) error {
	if err := r.pgRepo.DeleteOrder(ctx, id); err != nil {
		return err
	}

	key := strjoin.Join(getOrderKey, strconv.FormatUint(id, 10))
	r.logger.Error(r.lc.Delete(key))
	r.logger.Error(r.rc.Delete(ctx, key))
	r.logger.Error(r.rc.CFDel(ctx, cuckooFilter, id))
	return nil
}
