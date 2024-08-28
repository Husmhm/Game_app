package userservice

import (
	"fmt"
	"gameApp/entity"
	"gameApp/param"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	// TODO - We should verify phone number by verification code
	//validate phone number

	// validate name

	// TODO - chevk the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return param.RegisterResponse{}, fmt.Errorf("password lentgth should be greater than 8")
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
		return param.RegisterResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	// return created user
	//var resp RegisterResponse
	//resp.User.ID = createdUser.ID
	//resp.User.PhoneNumber = createdUser.PhoneNumber
	//resp.User.Name = createdUser.Name

	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}
