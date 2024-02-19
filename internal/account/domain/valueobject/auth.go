package valueobject

import "github.com/golang-jwt/jwt/v5"

// AuthPayload is a payload contains access token of customer
type AuthPayload struct {
	AccessToken string
}

// AuthResponse is a response contains customer id and expired status of access token
type AuthResponse struct {
	CustomerID uint64
	Expired    bool
}

type JWTClaims struct {
	CustomerID uint64
	Refresh    bool
	jwt.RegisteredClaims
}

// CustomerCredentials contains customer credentials
type CustomerCredentials struct {
	CustomerID uint64
	Active     bool
	Password   string
}
