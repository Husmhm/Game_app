package authservice

import (
	"fmt"
	"gameApp/entity"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessTokenSubject    string
	refreshTokenSubject   string
}

func New(signKey string, accessExpirationTime time.Duration, refreshExpirationTime time.Duration,
	accessTokenSubject string, refreshTokenSubject string) Service {
	return Service{
		signKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessTokenSubject:    accessTokenSubject,
		refreshTokenSubject:   refreshTokenSubject,
	}

}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.accessExpirationTime, s.accessTokenSubject)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.refreshExpirationTime, s.refreshTokenSubject)
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	fmt.Println("old string", bearerToken)
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)
	fmt.Println("new string", tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.signKey), nil
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
	tokenString, err := accessToken.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
