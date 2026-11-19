// Package handler is the Wails driving adapter for application-level
// information (version, updates) exposed to the frontend.
package handler

import (
	"context"
	"log"
	"time"

	"inventario/internal/modules/system/core/ports"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// UpdateProgressEvent is emitted while an update is being downloaded.
const UpdateProgressEvent = "update:progress"

type SystemHandler struct {
	svc ports.UpdateService
}

func NewSystemHandler(svc ports.UpdateService) *SystemHandler {
	return &SystemHandler{svc: svc}
}

// GetVersion returns the application version without the leading "v"
// (e.g. "1.2.3"), or "dev" for non-release builds.
func (h *SystemHandler) GetVersion() string {
	return h.svc.CurrentVersion()
}

// CheckForUpdate reports whether a newer release is published.
func (h *SystemHandler) CheckForUpdate(ctx context.Context) (UpdateInfoResponse, error) {
	info, err := h.svc.CheckForUpdate(ctx)
	if err != nil {
		return UpdateInfoResponse{}, err
	}
	return toUpdateInfoResponse(info), nil
}

// InstallUpdate downloads the installer for this platform, emitting
// UpdateProgressEvent, launches it and quits the application so the
// installer can replace it.
func (h *SystemHandler) InstallUpdate(ctx context.Context) error {
	err := h.svc.InstallUpdate(ctx, func(downloaded, total int64) {
		emit(UpdateProgressEvent, UpdateProgress{Downloaded: downloaded, Total: total})
	})
	if err != nil {
		return err
	}

	// Give the frontend time to receive the result before quitting.
	go func() {
		time.Sleep(time.Second)
		if app := application.Get(); app != nil {
			app.Quit()
		}
	}()
	return nil
}

// OpenReleasePage opens the latest release page in the default browser.
func (h *SystemHandler) OpenReleasePage() error {
	app := application.Get()
	if app == nil {
		return nil
	}
	return app.Browser.OpenURL(h.svc.ReleasePageURL())
}

func emit(name string, data any) {
	app := application.Get()
	if app == nil {
		return
	}
	if !app.Event.Emit(name, data) {
		log.Printf("evento %q no emitido", name)
	}
}
