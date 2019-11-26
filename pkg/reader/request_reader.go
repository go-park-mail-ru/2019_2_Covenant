package reader

import (
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

type ReqReader struct {
	v *validator.Validate
}

func NewReqReader() *ReqReader {
	return &ReqReader{
		v: validator.New(),
	}
}

func (rv *ReqReader) validate(req interface{}) error {
	err := rv.v.Struct(req)

	if err != nil {
		log.Print(err)
		return ErrBadParam
	}

	return nil
}

func (rv *ReqReader) Read(c echo.Context, request interface{}, check func(interface{}) bool) error {
	err := c.Bind(&request)

	if err != nil {
		return ErrUnprocessableEntity
	}

	if err := rv.validate(request); err != nil {
		return err
	}

	if check != nil && !check(request) {
		return ErrBadParam
	}

	return nil
}
