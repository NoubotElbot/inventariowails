package category

import "errors"

var (
	ErrNotFound      = errors.New("categoría no encontrada")
	ErrNameRequired  = errors.New("el nombre de la categoría es requerido")
	ErrNameDuplicate = errors.New("ya existe una categoría con ese nombre")
	ErrAlreadyDeleted = errors.New("la categoría ya está eliminada")
	ErrNotDeleted    = errors.New("la categoría no está eliminada")
)
