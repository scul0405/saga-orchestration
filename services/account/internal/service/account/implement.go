package account

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/internal/app"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
	"github.com/scul0405/saga-orchestration/services/account/internal/infrastructure/logger"
	"github.com/scul0405/saga-orchestration/services/account/internal/repository/postgres_repo"
)

type customerServiceImpl struct {
	repo   domain.CustomerRepository
	logger logger.Logger
}

// NewCustomerService returns new CustomerService
func NewCustomerService(repo domain.CustomerRepository, logger logger.Logger) app.CustomerService {
	return &customerServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *customerServiceImpl) GetPersonalInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerPersonalInfo, error) {
	info, err := s.repo.GetCustomerPersonalInfo(ctx, customerID)
	if err != nil {
		if err != postgres_repo.ErrCustomerNotFound {
			s.logger.Infof("GetPersonalInfo: failed to get personal info, err: %v", err)
		}

		return nil, err
	}

	return info, nil
}

func (s *customerServiceImpl) GetDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerDeliveryInfo, error) {
	info, err := s.repo.GetCustomerDeliveryInfo(ctx, customerID)
	if err != nil {
		if err != postgres_repo.ErrCustomerNotFound {
			s.logger.Infof("GetDeliveryInfo: failed to get delivery info, err: %v", err)
		}

		return nil, err
	}

	return info, nil
}

func (s *customerServiceImpl) UpdatePersonalInfo(ctx context.Context, customerID uint64, info *valueobject.CustomerPersonalInfo) error {
	return s.repo.UpdateCustomerPersonalInfo(ctx, customerID, info)
}

func (s *customerServiceImpl) UpdateDeliveryInfo(ctx context.Context, customerID uint64, info *valueobject.CustomerDeliveryInfo) error {
	return s.repo.UpdateCustomerDeliveryInfo(ctx, customerID, info)
}
