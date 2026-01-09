package service

import "errors"

var (
	ErrInvalidOrder         = errors.New("invalid order")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrAlreadyUploadedSelf  = errors.New("order already uploaded by user")
	ErrAlreadyUploadedOther = errors.New("order uploaded by another user")
)
