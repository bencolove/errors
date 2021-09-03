package errors

import (
	stderrors "errors"
)

// drop-in for standard errors
var (
	Is     = stderrors.Is
	As     = stderrors.As
	Unwrap = stderrors.Unwrap
)

type Stackerror *stackerror

func New(message string) error {
	return Wrap(message, nil)
}
