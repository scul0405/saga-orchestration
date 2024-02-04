package dto

import "github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"

type CustomerPersonalInfo struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type CustomerDeliveryInfo struct {
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (c *CustomerPersonalInfo) ToDomainObject() *valueobject.CustomerPersonalInfo {
	return &valueobject.CustomerPersonalInfo{
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
	}
}

func (c *CustomerDeliveryInfo) ToDomainObject() *valueobject.CustomerDeliveryInfo {
	return &valueobject.CustomerDeliveryInfo{
		Address:     c.Address,
		PhoneNumber: c.PhoneNumber,
	}
}
