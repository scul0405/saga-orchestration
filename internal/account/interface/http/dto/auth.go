package dto

import (
	"github.com/scul0405/saga-orchestration/internal/account/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/account/domain/valueobject"
)

type RegisterCustomer struct {
	Password    string `json:"password" binding:"required,min=8,max=128"`
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginCustomer struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *RegisterCustomer) ToDomainObject() entity.Customer {
	return entity.Customer{
		Password: r.Password,
		PersonalInfo: &valueobject.CustomerPersonalInfo{
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Email:     r.Email,
		},
		DeliveryInfo: &valueobject.CustomerDeliveryInfo{
			Address:     r.Address,
			PhoneNumber: r.PhoneNumber,
		},
	}
}
