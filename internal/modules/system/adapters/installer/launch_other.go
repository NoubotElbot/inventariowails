//go:build !windows && !darwin && !linux

package installer

import "inventario/internal/modules/system/core/domain"

func supported() bool { return false }

func launch(string) error {
	return domain.ErrManualUpdate
}
