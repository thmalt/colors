package colors

import "errors"

var (
	ErrInvalidSpace = errors.New("invalid space")
	ErrUnknownSpace = errors.New("unknown space")

	ErrInvalidConversion = errors.New("colors: invalid color conversion")
)
