package matchingvalidator

import (
	"fmt"
	"gameApp/entity"
	"gameApp/param"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddToWatingListRequest(req param.AddToWatingListRequest) (map[string]string, error) {
	const op = "matchingValidator.AddToWatingListRequest"
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category, validation.Required,
			validation.By(v.isCategoryValid)),
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

func (v Validator) isCategoryValid(value interface{}) error {
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrorMsgCategoryNotValid)
	}
	return nil

}
