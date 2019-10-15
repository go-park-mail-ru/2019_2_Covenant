package models

import "errors"

var (
	ErrInternalServerError  = 	errors.New("error: Internal Server Error")
	ErrBadParam 		    = 	errors.New("error: Bad Params")
	ErrNotFound 		    = 	errors.New("error: Not Found")
	ErrAlreadyExist 	    = 	errors.New("error: Already Exist")
)
