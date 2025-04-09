package apperror

import "errors"

var (
	ErrZipCodeRequired = errors.New("zipcode is required")
	ErrZipCodeInvalid  = errors.New("invalid zipcode")
	ErrZipCodeNotFound = errors.New("can not find zipcode")
)
