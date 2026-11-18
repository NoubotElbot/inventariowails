package domain

import "errors"

var (
	ErrProductRequired     = errors.New("el producto es requerido")
	ErrQuantityNotPositive = errors.New("la cantidad debe ser mayor a cero")
	ErrInvalidMovement     = errors.New("tipo de movimiento inválido")
)
