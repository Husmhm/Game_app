package userservice

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/phonenumber"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserById(id uint) (entity.User, error)
}

type Service struct {
	signkey string
	repo    Repository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func NewService(repo Repository, signkey string) Service {
	return Service{repo: repo, signkey: signkey}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - We should verify phone number by verification code
	//validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number isnt valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {

		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error : %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("Phone number is not unique")

		}

	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name isnt valid")
	}

	// TODO - chevk the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password lentgth should be greater than 8")
	}
	// hash password by bycript
	pass := []byte(req.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	hashStr := string(hash)
	if err != nil {
		panic(err)
	}
	// created new user in storage
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    hashStr,
	}
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	// return created user
	return RegisterResponse{
		User: createdUser,
	}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// generate jwt token
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// check existence phone number from repository
	// get the user by phone number
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("user does not exist")
	}
	// compare user.Password with the req.Password

	hErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if hErr != nil {
		return LoginResponse{}, fmt.Errorf("invalid password")
	}
	// return ok
	token, err := createToken(user.ID, s.signkey)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}
	return LoginResponse{AccessToken: token}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

// all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	// getUserByID
	user, err := s.repo.GetUserById(req.UserID)
	if err != nil {
		// I don't expect the repository call return "record not found" error,
		// because I assume the interactor input is sanitized.
		// TODO - we can use Rich Error.
		return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return ProfileResponse{Name: user.Name}, nil
}

type Claims struct {
	RegisteredClaims jwt.RegisteredClaims
	UserID           uint
}

func (c Claims) Valid() error {
	return nil
}

func createToken(userID uint, signKey string) (string, error) {
	// create a signer for rsa 256
	// TODO - replace with rsa 256 RS256 - https://github.com/golang-jwt/jwt/blob/main/http_example_test.go

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
