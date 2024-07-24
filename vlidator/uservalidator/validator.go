package uservalidator

import (
	"gameApp/dto"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/phonenumber"
	"gameApp/pkg/richerror"
)

type Repsitory interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repsitory
}

func New(repo Repsitory) Validator {
	return Validator{repo}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) error {
	const op = "userValidator.ValidateRegisterRequest"
	if !phonenumber.IsValid(req.PhoneNumber) {
		return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotValid).
			WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
		//return dto.RegisterResponse{}, fmt.Errorf("phone number isnt valid")
	}

	// check uniqueness of phone number
	if isUnique, err := v.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {

		if err != nil {
			return richerror.New(op).WithErr(err)
		}

		if !isUnique {
			return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotUnique).
				WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})

		}

		if len(req.Name) < 3 {
			return richerror.New(op).WithMessage(errmsg.ErrorMsgNameLength).WithKind(richerror.KindInvalid).
				WithMeta(map[string]interface{}{"name": req.Name})

		}

	}

}
