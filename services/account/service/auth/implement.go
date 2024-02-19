package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/scul0405/saga-orchestration/cmd/account/config"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/pkg/utils"
	"github.com/scul0405/saga-orchestration/services/account/app"
	"github.com/scul0405/saga-orchestration/services/account/domain"
	"github.com/scul0405/saga-orchestration/services/account/domain/entity"
	"github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
	"time"
)

var (
	ErrInvalidToken         = fmt.Errorf("invalid token")
	ErrCustomerNotFound     = fmt.Errorf("customer not found")
	ErrCustomerInactive     = fmt.Errorf("customer inactive")
	ErrTokenExpired         = fmt.Errorf("token expired")
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
)

type jwtAuthServiceImpl struct {
	jwtConfig config.JWTConfig
	repo      domain.JWTAuthRepository
	logger    logger.Logger
	sf        sonyflake.IDGenerator
}

// NewJWTAuthService returns a new instance of JWTAuthService
func NewJWTAuthService(jwtConfig config.JWTConfig, repo domain.JWTAuthRepository, logger logger.Logger, sf sonyflake.IDGenerator) app.AuthService {
	return &jwtAuthServiceImpl{
		jwtConfig: jwtConfig,
		repo:      repo,
		logger:    logger,
		sf:        sf,
	}
}
func (s *jwtAuthServiceImpl) Auth(ctx context.Context, authPayload *valueobject.AuthPayload) (*valueobject.AuthResponse, error) {
	token, err := s.parseToken(authPayload.AccessToken)

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return &valueobject.AuthResponse{
				Expired: true,
			}, nil
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*valueobject.JWTClaims)
	if !(ok && token.Valid) {
		return nil, ErrInvalidToken
	}

	if claims.Refresh {
		return nil, ErrInvalidToken
	}

	return &valueobject.AuthResponse{
		CustomerID: claims.CustomerID,
		Expired:    false,
	}, nil
}

func (s *jwtAuthServiceImpl) Register(ctx context.Context, customer *entity.Customer) (string, string, error) {
	sfID, err := s.sf.NextID()
	if err != nil {
		s.logger.Errorf("register: failed to generate customer id: %v", err)
		return "", "", err
	}

	customer.ID = sfID
	customer.Active = true

	if err = s.repo.CreateCustomer(ctx, customer); err != nil {
		s.logger.Errorf("register: failed to create customer: %v", err)
		return "", "", err
	}

	return s.generatePairToken(customer.ID)
}

func (s *jwtAuthServiceImpl) Login(ctx context.Context, email, password string) (string, string, error) {
	exist, customer, err := s.repo.GetCustomerCredentials(ctx, email)
	if err != nil {
		s.logger.Errorf("login: failed to get customer by email: %v", err)
		return "", "", err
	}

	if !exist {
		return "", "", ErrCustomerNotFound
	}

	if !customer.Active {
		return "", "", ErrCustomerInactive
	}

	if !utils.CheckPasswordHash(password, customer.Password) {
		return "", "", ErrAuthenticationFailed
	}

	return s.generatePairToken(customer.CustomerID)
}

func (s *jwtAuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := s.parseToken(refreshToken)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return "", "", ErrTokenExpired
		}

		return "", "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*valueobject.JWTClaims)
	if !(ok && token.Valid) {
		return "", "", ErrInvalidToken
	}

	if !claims.Refresh {
		return "", "", ErrInvalidToken
	}

	if !claims.ExpiresAt.After(time.Now()) {
		return "", "", ErrTokenExpired
	}

	exist, active, err := s.repo.CheckCustomer(ctx, claims.CustomerID)
	if err != nil {
		s.logger.Errorf("refresh token: failed to check customer: %v", err)
		return "", "", err
	}

	if !exist {
		return "", "", ErrCustomerNotFound
	}

	if !active {
		return "", "", ErrCustomerInactive
	}

	return s.generatePairToken(claims.CustomerID)
}

func (s *jwtAuthServiceImpl) parseToken(accessToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(accessToken, &valueobject.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.SecretKey), nil
	})
}

func (s *jwtAuthServiceImpl) generateToken(customerID uint64, refresh bool) (string, error) {
	var expiresAt time.Time
	if refresh {
		expiresAt = time.Now().Add(time.Duration(s.jwtConfig.RefreshTokenExpire) * time.Minute)
	} else {
		expiresAt = time.Now().Add(time.Duration(s.jwtConfig.AccessTokenExpire) * time.Minute)
	}

	claims := &valueobject.JWTClaims{
		CustomerID: customerID,
		Refresh:    refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.SecretKey))
}

func (s *jwtAuthServiceImpl) generatePairToken(customerID uint64) (string, string, error) {
	accessToken, err := s.generateToken(customerID, false)
	if err != nil {
		s.logger.Errorf("failed to generate access token: %v", err)
		return "", "", err
	}

	refreshToken, err := s.generateToken(customerID, true)
	if err != nil {
		s.logger.Errorf("failed to generate refresh token: %v", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
