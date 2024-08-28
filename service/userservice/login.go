package userservice

import (
	"fmt"
	"gameApp/param"
	"gameApp/pkg/richerror"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	hErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if hErr != nil {
		return param.LoginResponse{}, fmt.Errorf("invalid password")
	}
	// return ok
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return param.LoginResponse{User: param.UserInfo{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
	}, Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}}, nil

}
