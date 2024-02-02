package postgres_repo

import (
	"context"
	"errors"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
	"github.com/scul0405/saga-orchestration/services/account/internal/infrastructure/db/postgres/model"
	"github.com/scul0405/saga-orchestration/services/account/pkg"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEntry = errors.New("duplicate entry")
)

type customerStatus struct {
	Active bool
}

type CustomerCredentials struct {
	ID       uint64
	Active   bool
	Password string
}

type jwtAuthRepositoryImpl struct {
	db *gorm.DB
}

func NewJwtAuthRepositoryImpl(db *gorm.DB) domain.JWTAuthRepository {
	return &jwtAuthRepositoryImpl{
		db: db,
	}
}

func (r *jwtAuthRepositoryImpl) CheckCustomer(ctx context.Context, customerID uint64) (bool, bool, error) {
	var status customerStatus
	if err := r.db.Model(&model.Account{}).Select("active").
		Where("id = ?", customerID).First(&status).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}
	return true, status.Active, nil
}

func (r *jwtAuthRepositoryImpl) CreateCustomer(ctx context.Context, customer *entity.Customer) error {
	hashedPassword, err := pkg.HashPassword(customer.Password)
	if err != nil {
		return err
	}
	if err := r.db.Create(&model.Account{
		ID:          customer.ID,
		Active:      customer.Active,
		FirstName:   customer.PersonalInfo.FirstName,
		LastName:    customer.PersonalInfo.LastName,
		Email:       customer.PersonalInfo.Email,
		Address:     customer.DeliveryInfo.Address,
		PhoneNumber: customer.DeliveryInfo.PhoneNumber,
		Password:    hashedPassword,
	}).WithContext(ctx).Error; err != nil {
		return err
	}
	return nil
}

func (r *jwtAuthRepositoryImpl) GetCustomerCredentials(ctx context.Context, email string) (bool, *valueobject.CustomerCredentials, error) {
	var credentials CustomerCredentials
	if err := r.db.Model(&model.Account{}).Select("id", "active", "password").
		Where("email = ?", email).First(&credentials).WithContext(ctx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &valueobject.CustomerCredentials{
		CustomerID: credentials.ID,
		Active:     credentials.Active,
		Password:   credentials.Password,
	}, nil
}
