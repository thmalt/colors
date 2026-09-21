package colors

const (
	maxUint8  = 1<<8 - 1
	maxUint16 = 1<<16 - 1

	invMaxUint8  = 1.0 / maxUint8
	invMaxUint16 = 1.0 / maxUint16

	scale8To16 = 1<<8 + 1
)
