package utils

import (
	"errors"
	"time"

	"stock-trading/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(pUserID uuid.UUID, pRole string) (string, string, error) {
	lSecret := []byte(config.App.JWT.Secret)
	lExpHours := time.Duration(config.App.JWT.ExpirationHours) * time.Hour

	// Access Token
	lAccessClaims := JWTClaims{
		UserID: pUserID.String(),
		Role:   pRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(lExpHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "stock-trading-app",
			Subject:   pUserID.String(),
		},
	}
	lAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, lAccessClaims)
	lAccessTokenString, lErr := lAccessToken.SignedString(lSecret)
	if lErr != nil {
		return "", "", lErr
	}

	// Refresh Token (Longer lifespan)
	lRefreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(lExpHours * 24 * 7)), // 7 days
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   pUserID.String(),
	}
	lRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, lRefreshClaims)
	lRefreshTokenString, lErr := lRefreshToken.SignedString(lSecret)
	if lErr != nil {
		return "", "", lErr
	}

	return lAccessTokenString, lRefreshTokenString, nil
}

func ValidateToken(pTokenString string) (*JWTClaims, error) {
	lSecret := []byte(config.App.JWT.Secret)
	lToken, lErr := jwt.ParseWithClaims(pTokenString, &JWTClaims{}, func(pToken *jwt.Token) (interface{}, error) {
		return lSecret, nil
	})

	if lErr != nil {
		return nil, lErr
	}

	if lClaims, lOk := lToken.Claims.(*JWTClaims); lOk && lToken.Valid {
		return lClaims, nil
	}
	return nil, errors.New("invalid token")
}
