package authservice

import (
	"fmt"
	"gameApp/entity"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessTokenSubject    string        `koanf:"access_subject"`
	RefreshTokenSubject   string        `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg,
	}

}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessExpirationTime, s.config.AccessTokenSubject)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshExpirationTime, s.config.RefreshTokenSubject)
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	fmt.Println("old string", bearerToken)
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)
	fmt.Println("new string", tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) createToken(userID uint, expireDuratipn time.Duration, subject string) (string, error) {
	// create a signer for rsa 256
	// TODO - replace with rsa 256 RS256 - https://github.com/golang-jwt/jwt/blob/main/http_example_test.go

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuratipn)),
			Subject:   subject,
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
