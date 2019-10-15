package models

import "errors"

var (
	ErrInternalServerError  = 	errors.New("Error: Internal Server Error")
	ErrBadParam 		    = 	errors.New("Error: Bad Params")
	ErrNotFound 		    = 	errors.New("Error: Not Found")
	ErrAlreadyExist 	    = 	errors.New("Error: Already Exist")
)
