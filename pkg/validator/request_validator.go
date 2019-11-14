package validator

import (
	"2019_2_Covenant/internal/vars"
	"gopkg.in/go-playground/validator.v9"
)

type ReqValidator struct {
	v *validator.Validate
}

func NewReqValidator() *ReqValidator {
	return &ReqValidator{
		v: validator.New(),
	}
}

func (rv *ReqValidator) Validate(usr interface{}) error {
	err := rv.v.Struct(usr)

	if err != nil {
		return vars.ErrBadParam
	}

	return nil
}
