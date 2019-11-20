package validator

import (
	"2019_2_Covenant/internal/vars"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

type ReqReader struct {
	v *validator.Validate
}

func NewReqValidator() *ReqReader {
	return &ReqReader{
		v: validator.New(),
	}
}

func (rv *ReqReader) validate(req interface{}) error {
	err := rv.v.Struct(req)

	if err != nil {
		log.Print(err)
		return vars.ErrBadParam
	}

	return nil
}

func (rv *ReqReader) Read(c echo.Context, request interface{}) error {
	err := c.Bind(&request)

	if err != nil {
		return vars.ErrUnprocessableEntity
	}

	if err := rv.validate(request); err != nil {
		return err
	}

	return nil
}
