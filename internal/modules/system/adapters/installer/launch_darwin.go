package installer

import (
	"fmt"
	"os/exec"
)

func supported() bool { return true }

// launch mounts and opens the disk image in Finder so the user can drag the
// new app over the old one.
func launch(path string) error {
	if err := exec.Command("open", path).Start(); err != nil {
		return fmt.Errorf("abriendo imagen de disco: %w", err)
	}
	return nil
}
