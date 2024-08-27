package userservice

import (
	"fmt"
	"gameApp/dto"
	"gameApp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	hErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if hErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("invalid password")
	}
	// return ok
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return dto.LoginResponse{User: dto.UserInfo{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
	}, Tokens: dto.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}}, nil

}
