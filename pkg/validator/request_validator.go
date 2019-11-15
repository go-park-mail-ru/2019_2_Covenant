package validator

import (
	"2019_2_Covenant/internal/vars"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

type ReqValidator struct {
	v *validator.Validate
}

func NewReqValidator() *ReqValidator {
	return &ReqValidator{
		v: validator.New(),
	}
}

func (rv *ReqValidator) Validate(req interface{}) error {
	err := rv.v.Struct(req)

	if err != nil {
		log.Print(err)
		return vars.ErrBadParam
	}

	return nil
}
