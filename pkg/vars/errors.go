package vars

import "errors"

var (
	ErrInternalServerError  = 	errors.New("internal server error")
	ErrBadParam 		    = 	errors.New("bad params")
	ErrNotFound 		    = 	errors.New("not found")
	ErrAlreadyExist 	    = 	errors.New("already exist")
	ErrExpired              =   errors.New("session expired")
)
