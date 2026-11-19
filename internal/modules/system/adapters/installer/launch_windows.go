package installer

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func supported() bool { return true }

// launch opens the NSIS installer through ShellExecute so Windows shows the
// UAC prompt for per-machine installs (CreateProcess would fail with
// ERROR_ELEVATION_REQUIRED).
func launch(path string) error {
	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err := windows.ShellExecute(0, verb, file, nil, nil, windows.SW_SHOWNORMAL); err != nil {
		return fmt.Errorf("ejecutando instalador: %w", err)
	}
	return nil
}
