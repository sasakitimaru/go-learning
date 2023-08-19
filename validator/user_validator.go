package validator

import (
	"go-rest-api/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IUserValidator interface {
	UserValidate(user model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (uv *userValidator) UserValidate(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required.Error("email is required"), validation.Length(1, 30).Error("limited max 10 characters")),
		validation.Field(&user.Password, validation.Required.Error("password is required"), validation.Length(6, 30).Error("limited min 6 max 10 characters")),
	)
}
