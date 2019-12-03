package vars

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrBadParam            = errors.New("bad params")
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExist        = errors.New("already exist")
	ErrExpired             = errors.New("session expired")
	ErrRetrievingError     = errors.New("retrieving error")
	ErrBadCSRF             = errors.New("csrf error")
	ErrUnathorized         = errors.New("unauthorized")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrPermissionDenied     = errors.New("permission denied")
)
