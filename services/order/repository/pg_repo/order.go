package pg_repo

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/order/domain"
	"github.com/scul0405/saga-orchestration/services/order/domain/entity"
	"github.com/scul0405/saga-orchestration/services/order/domain/valueobject"
	"github.com/scul0405/saga-orchestration/services/order/infrastructure/db/postgres/model"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) GetOrder(ctx context.Context, id uint64) (*entity.Order, error) {
	var order []model.Order
	if err := r.db.Where("id = ?", id).Find(&order).WithContext(ctx).Error; err != nil {
		return nil, err
	}

	purchasedProducts := make([]valueobject.PurchasedProduct, len(order))

	for i, entry := range order {
		purchasedProducts[i] = valueobject.PurchasedProduct{
			ID:       entry.ProductID,
			Quantity: entry.Quantity,
		}
	}

	return &entity.Order{
		ID:                order[0].ID,
		CustomerID:        order[0].CustomerID,
		PurchasedProducts: &purchasedProducts,
	}, nil
}

func (r *orderRepositoryImpl) CreateOrder(ctx context.Context, order *entity.Order) error {
	entries := make([]model.Order, len(*(order.PurchasedProducts)))

	for i, product := range *(order.PurchasedProducts) {
		entries[i] = model.Order{
			ID:         order.ID,
			ProductID:  product.ID,
			CustomerID: order.CustomerID,
			Quantity:   product.Quantity,
		}
	}

	if err := r.db.Create(&entries).WithContext(ctx).Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepositoryImpl) DeleteOrder(ctx context.Context, id uint64) error {
	if err := r.db.Exec("DELETE FROM orders WHERE id = ?", id).WithContext(ctx).Error; err != nil {
		return err
	}
	return nil
}
