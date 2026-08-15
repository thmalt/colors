package gradient

import "errors"

var (
	ErrInvalidTransform = errors.New("invalid transform")
	ErrInvalidSize      = errors.New("invalid size")
	ErrInvalidRadius    = errors.New("invalid radius")
	ErrInvalidCenter    = errors.New("invalid center")
	ErrInvalidFocal     = errors.New("invalid focal point")
	ErrInvalidAngle     = errors.New("invalid angle")
)
