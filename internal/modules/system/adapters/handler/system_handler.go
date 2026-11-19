// Package handler is the Wails driving adapter for application-level
// information (version, build metadata) exposed to the frontend.
package handler

import (
	"strings"

	"inventario/internal/platform/buildinfo"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetVersion returns the application version without the leading "v"
// (e.g. "1.2.3"), or "dev" for non-release builds.
func (h *SystemHandler) GetVersion() string {
	return strings.TrimPrefix(buildinfo.Version, "v")
}
