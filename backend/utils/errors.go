package utils

import "errors"

var (
	ErrInvalidDateFormat       = errors.New("invalid date format, use YYYY-MM-DD")
	ErrQueueNotFound           = errors.New("queue not found")
	ErrCannotDelete            = errors.New("only waiting or cancelled queue can be deleted")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrServiceNotFound         = errors.New("service not found")
	ErrInvalidEstimatedMinutes = errors.New("estimated minutes must be greater than 0")
)
