package ports

import (
	"context"

	"inventario/internal/modules/system/core/domain"
)

// ProgressFunc reports download progress. total is -1 when unknown.
type ProgressFunc func(downloaded, total int64)

// ReleaseProvider is an outbound port that fetches the latest published
// release (e.g. from GitHub Releases).
type ReleaseProvider interface {
	LatestRelease(ctx context.Context) (domain.Release, error)
}

// Installer is an outbound port that downloads an installer asset and hands
// it off to the operating system. On success the caller is expected to quit
// the application so the installer can replace it.
type Installer interface {
	Install(ctx context.Context, asset domain.Asset, progress ProgressFunc) error
}

// UpdateService is the inbound (driving) port for version and update use cases.
type UpdateService interface {
	CurrentVersion() string
	CheckForUpdate(ctx context.Context) (domain.UpdateInfo, error)
	InstallUpdate(ctx context.Context, progress ProgressFunc) error
	ReleasePageURL() string
}
