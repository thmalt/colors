package gradient

import "errors"

var (
	ErrInvalidTransform = errors.New("gradient: invalid transform")
	ErrInvalidSize      = errors.New("gradient: invalid size")
	ErrInvalidRadius    = errors.New("gradient: invalid radius")
	ErrInvalidCenter    = errors.New("gradient: invalid center")
	ErrInvalidFocal     = errors.New("gradient: invalid focal point")
	ErrInvalidAngle     = errors.New("gradient: invalid angle")
)
