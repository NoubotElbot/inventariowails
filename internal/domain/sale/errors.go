package sale

import "errors"

var (
	ErrNotFound        = errors.New("venta no encontrada")
	ErrNoItems         = errors.New("la venta debe tener al menos un producto")
	ErrQuantityInvalid = errors.New("la cantidad debe ser mayor a cero")
	ErrProductInvalid  = errors.New("producto inválido en la venta")
)
