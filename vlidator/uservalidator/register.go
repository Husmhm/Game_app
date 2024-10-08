package uservalidator

import (
	"fmt"
	"gameApp/param"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (map[string]string, error) {
	const op = "userValidator.ValidateRegisterRequest"
	if err := validation.ValidateStruct(&req,
		// TODO - add 3 to config
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`))),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.checkPhoneNumberUniqueness)),
	); err != nil {
		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)
	// check uniqueness of phone number
	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {

		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)

		}

	}
	return nil
}
