package stock

import "errors"

var (
	ErrProductRequired = errors.New("el producto es requerido")
	ErrQuantityZero    = errors.New("la cantidad no puede ser cero")
	ErrInvalidMovement = errors.New("tipo de movimiento inválido")
)
