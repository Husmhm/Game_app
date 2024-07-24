package userservice

import (
	"fmt"
	"gameApp/dto"
	"gameApp/entity"
	"gameApp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserById(id uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo Repository
}

func NewService(authGenerator AuthGenerator, repo Repository) Service {
	return Service{auth: authGenerator, repo: repo}
}

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO - We should verify phone number by verification code
	//validate phone number

	// validate name

	// TODO - chevk the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return dto.RegisterResponse{}, fmt.Errorf("password lentgth should be greater than 8")
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
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	// return created user
	//var resp RegisterResponse
	//resp.User.ID = createdUser.ID
	//resp.User.PhoneNumber = createdUser.PhoneNumber
	//resp.User.Name = createdUser.Name

	return dto.RegisterResponse{User: dto.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// generate jwt token
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type LoginResponse struct {
	User   dto.UserInfo `json:"user"`
	Tokens Tokens       `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
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
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return LoginResponse{User: dto.UserInfo{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
	}, Tokens: Tokens{AccessToken: accessToken, RefreshToken: refreshToken}}, nil

}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

// all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	const op = "userService.Profile"
	user, err := s.repo.GetUserById(req.UserID)
	if err != nil {

		return ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

	}

	return ProfileResponse{Name: user.Name}, nil
}
