package domain

import "errors"

var (
	ErrNoReleases    = errors.New("no hay versiones publicadas")
	ErrNoUpdate      = errors.New("no hay actualizaciones disponibles")
	ErrManualUpdate  = errors.New("esta instalación no admite actualización automática; descarga la nueva versión manualmente")
	ErrUpdateRunning = errors.New("ya hay una actualización en curso")
	ErrChecksum      = errors.New("la descarga no coincide con el checksum publicado; la actualización se canceló")
)
