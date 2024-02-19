package account

import (
	"context"
	"github.com/scul0405/saga-orchestration/cmd/account/config"
	"github.com/scul0405/saga-orchestration/pkg/appconfig"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
	"github.com/scul0405/saga-orchestration/services/account/repository/postgres_repo"
	"github.com/scul0405/saga-orchestration/services/account/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetPersonalInfo(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	testcases := []struct {
		name      string
		buildStub func(repo *mock.MockCustomerRepository)
	}{
		{
			name: "success",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetCustomerPersonalInfo(gomock.Any(), gomock.Eq(customerID)).Return(&valueobject.CustomerPersonalInfo{}, nil)
			},
		},
		{
			name: "not found",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetCustomerPersonalInfo(gomock.Any(), gomock.Eq(customerID)).Return(nil, postgres_repo.ErrCustomerNotFound)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := &config.Config{
				App: appconfig.App{
					Logger: appconfig.Logger{
						Development:       true,
						DisableCaller:     false,
						DisableStacktrace: false,
						Encoding:          "json",
					},
				},
			}
			apiLogger := logger.NewApiLogger(&cfg.App)

			repo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStub(repo)
			service := NewCustomerService(repo, apiLogger)

			ctx := context.Background()
			info, err := service.GetPersonalInfo(ctx, customerID)
			if err != nil {
				require.Equal(t, postgres_repo.ErrCustomerNotFound, err)
			} else {
				require.NotNil(t, info)
			}
		})
	}
}

func TestGetDeliveryInfo(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	testcases := []struct {
		name      string
		buildStub func(repo *mock.MockCustomerRepository)
	}{
		{
			name: "success",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetCustomerDeliveryInfo(gomock.Any(), gomock.Eq(customerID)).Return(&valueobject.CustomerDeliveryInfo{}, nil)
			},
		},
		{
			name: "not found",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().GetCustomerDeliveryInfo(gomock.Any(), gomock.Eq(customerID)).Return(nil, postgres_repo.ErrCustomerNotFound)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := &config.Config{
				App: appconfig.App{
					Logger: appconfig.Logger{
						Development:       true,
						DisableCaller:     false,
						DisableStacktrace: false,
						Encoding:          "json",
					},
				},
			}
			apiLogger := logger.NewApiLogger(&cfg.App)

			repo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStub(repo)
			service := NewCustomerService(repo, apiLogger)

			ctx := context.Background()
			info, err := service.GetDeliveryInfo(ctx, customerID)
			if err != nil {
				require.Equal(t, postgres_repo.ErrCustomerNotFound, err)
			} else {
				require.NotNil(t, info)
			}
		})
	}
}

func TestUpdatePersonalInfo(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	infoBase := &valueobject.CustomerPersonalInfo{
		FirstName: "dep",
		LastName:  "trai",
		Email:     "deptrai@gmail.com",
	}

	testcases := []struct {
		name      string
		buildStub func(repo *mock.MockCustomerRepository)
	}{
		{
			name: "success",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().UpdateCustomerPersonalInfo(gomock.Any(), gomock.Eq(customerID), gomock.Eq(infoBase)).Return(nil)
			},
		},
		{
			name: "not found",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().UpdateCustomerPersonalInfo(gomock.Any(), gomock.Eq(customerID), gomock.Eq(infoBase)).Return(postgres_repo.ErrCustomerNotFound)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := &config.Config{
				App: appconfig.App{
					Logger: appconfig.Logger{
						Development:       true,
						DisableCaller:     false,
						DisableStacktrace: false,
						Encoding:          "json",
					},
				},
			}
			apiLogger := logger.NewApiLogger(&cfg.App)

			repo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStub(repo)
			service := NewCustomerService(repo, apiLogger)

			ctx := context.Background()
			err := service.UpdatePersonalInfo(ctx, customerID, infoBase)
			if err != nil {
				require.Equal(t, postgres_repo.ErrCustomerNotFound, err)
			}
		})
	}
}

func TestUpdateDeliveryInfo(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	infoBase := &valueobject.CustomerDeliveryInfo{
		Address:     "123 abc",
		PhoneNumber: "123456789",
	}

	testcases := []struct {
		name      string
		buildStub func(repo *mock.MockCustomerRepository)
	}{
		{
			name: "success",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().UpdateCustomerDeliveryInfo(gomock.Any(), gomock.Eq(customerID), gomock.Eq(infoBase)).Return(nil)
			},
		},
		{
			name: "not found",
			buildStub: func(repo *mock.MockCustomerRepository) {
				repo.EXPECT().UpdateCustomerDeliveryInfo(gomock.Any(), gomock.Eq(customerID), gomock.Eq(infoBase)).Return(postgres_repo.ErrCustomerNotFound)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := &config.Config{
				App: appconfig.App{
					Logger: appconfig.Logger{
						Development:       true,
						DisableCaller:     false,
						DisableStacktrace: false,
						Encoding:          "json",
					},
				},
			}
			apiLogger := logger.NewApiLogger(&cfg.App)

			repo := mock.NewMockCustomerRepository(ctrl)
			tc.buildStub(repo)
			service := NewCustomerService(repo, apiLogger)

			ctx := context.Background()
			err := service.UpdateDeliveryInfo(ctx, customerID, infoBase)
			if err != nil {
				require.Equal(t, postgres_repo.ErrCustomerNotFound, err)
			}
		})
	}
}
