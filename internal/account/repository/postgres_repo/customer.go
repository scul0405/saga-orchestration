package postgres_repo

import (
	"context"
	"errors"
	"github.com/scul0405/saga-orchestration/internal/account/domain"
	"github.com/scul0405/saga-orchestration/internal/account/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/account/infrastructure/db/postgres/model"
	"gorm.io/gorm"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

// CustomerRepositoryImpl implements CustomerRepository interface
type customerRepositoryImpl struct {
	db *gorm.DB
}

// NewCustomerRepositoryImpl returns new CustomerRepositoryImpl
func NewCustomerRepositoryImpl(db *gorm.DB) domain.CustomerRepository {
	return &customerRepositoryImpl{
		db: db,
	}
}

func (r *customerRepositoryImpl) GetCustomerPersonalInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerPersonalInfo, error) {
	var info valueobject.CustomerPersonalInfo
	if err := r.db.Model(&model.Account{}).Select("first_name", "last_name", "email").
		Where("id = ? AND active = TRUE", customerID).First(&info).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &info, nil
}

func (r *customerRepositoryImpl) GetCustomerDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerDeliveryInfo, error) {
	var info valueobject.CustomerDeliveryInfo
	if err := r.db.Model(&model.Account{}).Select("address", "phone_number").
		Where("id = ? AND active = TRUE", customerID).First(&info).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &info, nil
}

func (r *customerRepositoryImpl) UpdateCustomerPersonalInfo(ctx context.Context, customerID uint64, personalInfo *valueobject.CustomerPersonalInfo) error {
	if err := r.db.Model(&model.Account{}).Where("id = ? AND active = TRUE", customerID).
		Updates(model.Account{
			FirstName: personalInfo.FirstName,
			LastName:  personalInfo.LastName,
			Email:     personalInfo.Email,
		}).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomerNotFound
		}
		return err
	}
	return nil
}

func (r *customerRepositoryImpl) UpdateCustomerDeliveryInfo(ctx context.Context, customerID uint64, deliveryInfo *valueobject.CustomerDeliveryInfo) error {
	if err := r.db.Model(&model.Account{}).Where("id = ? AND active = TRUE", customerID).
		Updates(model.Account{
			Address:     deliveryInfo.Address,
			PhoneNumber: deliveryInfo.PhoneNumber,
		}).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomerNotFound
		}
		return err
	}
	return nil
}
