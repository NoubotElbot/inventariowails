package sqlite

import (
	"strings"
)

// isUniqueConstraint detecta violaciones de UNIQUE en SQLite.
// SQLite no expone códigos de error tipados via database/sql,
// por lo que se inspecciona el mensaje de error.
func isUniqueConstraint(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed")
}
