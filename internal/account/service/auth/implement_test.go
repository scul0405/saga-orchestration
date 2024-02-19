package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scul0405/saga-orchestration/cmd/account/config"
	"github.com/scul0405/saga-orchestration/internal/account/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/account/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/account/service/mock"
	"github.com/scul0405/saga-orchestration/pkg/appconfig"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	cfg := &config.Config{
		App: appconfig.App{
			Logger: appconfig.Logger{
				Development:       true,
				DisableCaller:     false,
				DisableStacktrace: false,
				Encoding:          "json",
			},
		},
		JWTConfig: config.JWTConfig{
			SecretKey:          "test-secret-key-for-jwt-12345678",
			AccessTokenExpire:  5,
			RefreshTokenExpire: 15,
		},
	}

	apiLogger := logger.NewApiLogger(&cfg.App)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	testcases := []struct {
		name          string
		setupAuth     func() *valueobject.AuthPayload
		checkResponse func(t *testing.T, info *valueobject.AuthResponse, err error)
	}{
		{
			name: "success",
			setupAuth: func() *valueobject.AuthPayload {
				token, err := generateToken(customerID, false, 5, cfg.JWTConfig.SecretKey)
				require.NoError(t, err)
				return &valueobject.AuthPayload{
					AccessToken: token,
				}
			},
			checkResponse: func(t *testing.T, info *valueobject.AuthResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, info)

				require.False(t, info.Expired)
				require.Equal(t, customerID, info.CustomerID)
			},
		},
		{
			name: "token expired",
			setupAuth: func() *valueobject.AuthPayload {
				token, err := generateToken(customerID, false, -5, cfg.JWTConfig.SecretKey)
				require.NoError(t, err)
				return &valueobject.AuthPayload{
					AccessToken: token,
				}
			},
			checkResponse: func(t *testing.T, info *valueobject.AuthResponse, err error) {
				require.Equal(t, err, ErrInvalidToken)
			},
		},
		{
			name: "invalid token",
			setupAuth: func() *valueobject.AuthPayload {
				token := "invalid"
				return &valueobject.AuthPayload{
					AccessToken: token,
				}
			},
			checkResponse: func(t *testing.T, info *valueobject.AuthResponse, err error) {
				require.Equal(t, err, ErrInvalidToken)
			},
		},
		{
			name: "refresh token",
			setupAuth: func() *valueobject.AuthPayload {
				token, err := generateToken(customerID, true, 5, cfg.JWTConfig.SecretKey)
				require.NoError(t, err)
				return &valueobject.AuthPayload{
					AccessToken: token,
				}
			},
			checkResponse: func(t *testing.T, info *valueobject.AuthResponse, err error) {
				require.Equal(t, err, ErrInvalidToken)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockJWTAuthRepository(ctrl)
			service := NewJWTAuthService(cfg.JWTConfig, repo, apiLogger, sf)
			ctx := context.Background()

			payload := tc.setupAuth()
			info, err := service.Auth(ctx, payload)
			tc.checkResponse(t, info, err)
		})
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	cfg := &config.Config{
		App: appconfig.App{
			Logger: appconfig.Logger{
				Development:       true,
				DisableCaller:     false,
				DisableStacktrace: false,
				Encoding:          "json",
			},
		},
		JWTConfig: config.JWTConfig{
			SecretKey:          "test-secret-key-for-jwt-12345678",
			AccessTokenExpire:  5,
			RefreshTokenExpire: 15,
		},
	}

	apiLogger := logger.NewApiLogger(&cfg.App)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	testcases := []struct {
		name          string
		buildStub     func(repo *mock.MockJWTAuthRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "success",
			buildStub: func(repo *mock.MockJWTAuthRepository) {
				repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockJWTAuthRepository(ctrl)
			service := NewJWTAuthService(cfg.JWTConfig, repo, apiLogger, sf)
			ctx := context.Background()

			if tc.buildStub != nil {
				tc.buildStub(repo)
			}

			customer := &entity.Customer{
				ID:       customerID,
				Active:   true,
				Password: "secret",
				PersonalInfo: &valueobject.CustomerPersonalInfo{
					FirstName: "dep",
					LastName:  "trai",
					Email:     "deptrai@gmail.com",
				},
				DeliveryInfo: &valueobject.CustomerDeliveryInfo{
					Address:     "123 abc",
					PhoneNumber: "123456789",
				},
			}

			_, _, err = service.Register(ctx, customer)
			tc.checkResponse(t, err)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	cfg := &config.Config{
		App: appconfig.App{
			Logger: appconfig.Logger{
				Development:       true,
				DisableCaller:     false,
				DisableStacktrace: false,
				Encoding:          "json",
			},
		},
		JWTConfig: config.JWTConfig{
			SecretKey:          "test-secret-key-for-jwt-12345678",
			AccessTokenExpire:  5,
			RefreshTokenExpire: 15,
		},
	}

	apiLogger := logger.NewApiLogger(&cfg.App)

	customerID, err := sf.NextID()
	require.NoError(t, err)

	type LoginData struct {
		Email    string
		Password string
	}

	password := "secret"

	hashPw, err := utils.HashPassword(password)
	require.NoError(t, err)
	customerCred := &valueobject.CustomerCredentials{
		CustomerID: customerID,
		Active:     true,
		Password:   hashPw,
	}

	testcases := []struct {
		name           string
		setupLoginData func() LoginData
		buildStub      func(repo *mock.MockJWTAuthRepository)
		checkResponse  func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupLoginData: func() LoginData {
				return LoginData{
					Email:    "deptrai@gmail.com",
					Password: password,
				}
			},
			buildStub: func(repo *mock.MockJWTAuthRepository) {
				repo.EXPECT().GetCustomerCredentials(gomock.Any(), gomock.Any()).
					Return(true, customerCred, nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockJWTAuthRepository(ctrl)
			service := NewJWTAuthService(cfg.JWTConfig, repo, apiLogger, sf)
			ctx := context.Background()

			loginData := tc.setupLoginData()

			if tc.buildStub != nil {
				tc.buildStub(repo)
			}

			_, _, err = service.Login(ctx, loginData.Email, loginData.Password)
			tc.checkResponse(t, err)
		})
	}
}

func generateToken(customerID uint64, refresh bool, expire int64, secretKey string) (string, error) {
	claims := &valueobject.JWTClaims{
		CustomerID: customerID,
		Refresh:    refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
