package uservalidator

import "gameApp/entity"

const (
	phoneNumberRegex = "^09[0-9]{9}$"
)

type Repsitory interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repsitory
}

func New(repo Repsitory) Validator {
	return Validator{repo}
}
