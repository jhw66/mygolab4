package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenType    = "access"
	RefreshTokenType   = "refresh"
	ChallengeTokenType = "challenge"
)

type myClaim struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, expireSeconds int64, userID string, tokenType string) (string, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", errors.New("secret cannot be empty")
	}
	if expireSeconds <= 0 {
		return "", errors.New("expireSeconds must be positive")
	}

	claims := myClaim{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString, expectedTokenType string) (*myClaim, error) {
	secret = strings.TrimSpace(secret)
	tokenString = strings.TrimSpace(tokenString)
	if secret == "" || tokenString == "" {
		return nil, errors.New("secret and tokenString cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &myClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*myClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedTokenType {
		return nil, errors.New("unexpected token type")
	}

	return claims, nil
}

func GetTokenFromHeader(r *http.Request) (tokenString string, err error) {
	bearerToken := strings.TrimSpace(r.Header.Get("Authorization"))
	if bearerToken == "" {
		return "", errors.New("token not found in header")
	}

	parts := strings.SplitN(bearerToken, " ", 2)
	if len(parts) != 2 || strings.EqualFold(parts[0], "Bearer") == false {
		return "", errors.New("invalid token format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
