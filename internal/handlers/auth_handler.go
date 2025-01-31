package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mahabub618/minipack/config"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type SignUpRequest struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Phone          string `json:"phone"`
	AddressStreet  string `json:"address_street"`
	AddressCity    string `json:"address_city"`
	AddressCountry string `json:"address_country"`
}

type SignupResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone,omitempty"`
	AddressStreet  string `json:"address_street,omitempty"`
	AddressCity    string `json:"address_city,omitempty"`
	AddressCountry string `json:"address_country,omitempty"`
	Role           string `json:"role"`
}

func GenerateAccessToken(userID int, role string, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(cfg.AccessTokenExpiry)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func GenerateRefreshToken(userID int, cfg *config.Config) (string, time.Time, error) {
	expirationTime := time.Now().Add(cfg.RefreshTokenExpiry)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = expirationTime.Unix()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	return tokenString, expirationTime, err
}
