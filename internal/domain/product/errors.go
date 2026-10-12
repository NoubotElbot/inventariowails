package product

import "errors"

var (
	ErrNotFound          = errors.New("producto no encontrado")
	ErrNameRequired      = errors.New("el nombre del producto es requerido")
	ErrCodeRequired      = errors.New("el código del producto es requerido")
	ErrCodeDuplicate     = errors.New("ya existe un producto con ese código")
	ErrPriceInvalid      = errors.New("el precio debe ser mayor a cero")
	ErrStockNegative     = errors.New("el stock no puede ser negativo")
	ErrInsufficientStock = errors.New("stock insuficiente")
)
