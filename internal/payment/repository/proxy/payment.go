package proxy

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/payment/domain"
	"github.com/scul0405/saga-orchestration/internal/payment/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/payment/repository/pgrepo"
	"github.com/scul0405/saga-orchestration/internal/pkg/cache"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/strjoin"
	"gorm.io/gorm"
	"strconv"
)

type paymentRepositoryImpl struct {
	pgRepo pgrepo.PaymentRepository
	lc     cache.LocalCache
	rc     cache.RedisCache
	logger logger.Logger
}

func NewPaymentRepository(
	pgRepo pgrepo.PaymentRepository,
	lc cache.LocalCache,
	rc cache.RedisCache,
	logger logger.Logger) (domain.PaymentRepository, error) {

	exist, err := rc.CFExist(context.Background(), cuckooFilter, dummnyItem)
	if err != nil {
		logger.Error("Payment: failed to check cuckoo filter existence", err)
		return nil, err
	}

	if !exist {
		if err = rc.CFReserve(context.Background(), cuckooFilter, 1000, 4, 100); err != nil {
			logger.Error("Payment: failed to reserve cuckoo filter", err)
			return nil, err
		}
		if err = rc.CFAdd(context.Background(), cuckooFilter, dummnyItem); err != nil {
			logger.Error("Payment: failed to add dummy item to cuckoo filter", err)
			return nil, err
		}

		logger.Info("Payment: cuckoo filter created")
	} else {
		logger.Info("Payment: cuckoo filter already exist")
	}

	return &paymentRepositoryImpl{
		pgRepo: pgRepo,
		lc:     lc,
		rc:     rc,
		logger: logger,
	}, nil
}

func (r *paymentRepositoryImpl) GetPayment(ctx context.Context, id uint64) (*entity.Payment, error) {
	payment := &entity.Payment{}
	key := strjoin.Join(getPaymentKey, strconv.FormatUint(id, 10))
	ok, err := r.lc.Get(key, payment)
	if ok && err == nil {
		return payment, nil
	}

	exist, err := r.rc.CFExist(ctx, cuckooFilter, id)
	r.logger.Error(err)
	if !exist && err == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ok, err = r.rc.Get(ctx, key, payment)
	r.logger.Error(err)
	if ok && err == nil {
		r.logger.Error(r.lc.Set(key, payment))
		return payment, nil
	}

	payment, err = r.pgRepo.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = r.lc.Set(key, payment); err != nil {
		r.logger.Error("Payment: failed to set payment to local cache", err)
	}

	r.logger.Error(r.rc.Set(ctx, key, payment))
	return payment, nil
}

func (r *paymentRepositoryImpl) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	err := r.pgRepo.CreatePayment(ctx, payment)
	if err != nil {
		return err
	}

	r.logger.Error(r.rc.CFAdd(ctx, cuckooFilter, payment.ID))
	return nil
}

func (r *paymentRepositoryImpl) DeletePayment(ctx context.Context, id uint64) error {
	if err := r.pgRepo.DeletePayment(ctx, id); err != nil {
		return err
	}

	key := strjoin.Join(getPaymentKey, strconv.FormatUint(id, 10))
	r.logger.Error(r.lc.Delete(key))
	r.logger.Error(r.rc.Delete(ctx, key))
	r.logger.Error(r.rc.CFDel(ctx, cuckooFilter, id))
	return nil
}
