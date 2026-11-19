package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"inventario/internal/modules/system/core/domain"
)

// supported reports whether the app runs as an AppImage, the only Linux
// format that can update itself. deb/rpm/arch installs are managed by the
// system package manager.
func supported() bool { return os.Getenv("APPIMAGE") != "" }

// launch replaces the running AppImage with the downloaded one and starts
// it.
func launch(path string) error {
	defer os.RemoveAll(filepath.Dir(path))

	target := os.Getenv("APPIMAGE")
	if target == "" {
		return domain.ErrManualUpdate
	}

	// Copy next to the target first so the final rename is atomic (os.TempDir
	// may live on a different filesystem).
	tmp := target + ".new"
	if err := copyFile(path, tmp, 0o755); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("preparando AppImage: %w", err)
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("reemplazando AppImage: %w", err)
	}

	if err := exec.Command(target).Start(); err != nil {
		return fmt.Errorf("iniciando nueva versión: %w", err)
	}
	return nil
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
