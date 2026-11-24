// Package services implements the system module use cases: reporting the
// running version and checking/installing updates from published releases.
package services

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"inventario/internal/modules/system/core/domain"
	"inventario/internal/modules/system/core/ports"
)

var _ ports.UpdateService = (*UpdateService)(nil)

type UpdateService struct {
	current        string
	goos, goarch   string
	releasesURL    string
	releases       ports.ReleaseProvider
	installer      ports.Installer
	installRunning atomic.Bool

	mu   sync.Mutex
	last *domain.UpdateInfo
}

// Option customizes an UpdateService.
type Option func(*UpdateService)

// WithPlatform overrides the target platform (defaults to the running one).
func WithPlatform(goos, goarch string) Option {
	return func(s *UpdateService) {
		s.goos, s.goarch = goos, goarch
	}
}

// NewUpdateService builds the service. current is the running version
// ("v1.2.3", "1.2.3" or "dev"); releasesURL is the fallback page opened when
// no specific release is known.
func NewUpdateService(
	current, releasesURL string,
	releases ports.ReleaseProvider,
	installer ports.Installer,
	opts ...Option,
) *UpdateService {
	s := &UpdateService{
		current:     current,
		goos:        runtime.GOOS,
		goarch:      runtime.GOARCH,
		releasesURL: releasesURL,
		releases:    releases,
		installer:   installer,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// CurrentVersion returns the running version without the leading "v"
// (e.g. "1.2.3"), or "dev" for non-release builds.
func (s *UpdateService) CurrentVersion() string {
	return strings.TrimPrefix(s.current, "v")
}

// CheckForUpdate queries the latest release and compares it with the running
// version. Non-release builds never report updates and make no requests.
func (s *UpdateService) CheckForUpdate(ctx context.Context) (domain.UpdateInfo, error) {
	info := domain.UpdateInfo{CurrentVersion: s.CurrentVersion()}

	current, err := domain.ParseVersion(s.current)
	if err != nil {
		return info, nil
	}

	release, err := s.releases.LatestRelease(ctx)
	if err != nil {
		return info, err
	}

	latest, err := domain.ParseVersion(release.Tag)
	if err != nil {
		return info, err
	}

	info.LatestVersion = latest.String()
	info.Notes = release.Notes
	info.ReleaseURL = release.URL
	info.PublishedAt = release.PublishedAt
	info.Available = latest.Compare(current) > 0
	if info.Available {
		// Without a published checksum the download cannot be verified, so
		// the update falls back to the manual path.
		if asset := domain.SelectAsset(release.Assets, s.goos, s.goarch); asset != nil && asset.SHA256 != "" {
			info.Asset = asset
		}
	}

	s.mu.Lock()
	s.last = &info
	s.mu.Unlock()

	return info, nil
}

// InstallUpdate downloads the installer of the newest release and launches
// it. It returns domain.ErrManualUpdate when the platform has no compatible
// installer, in which case the user should download it from ReleasePageURL.
func (s *UpdateService) InstallUpdate(ctx context.Context, progress ports.ProgressFunc) error {
	if !s.installRunning.CompareAndSwap(false, true) {
		return domain.ErrUpdateRunning
	}
	defer s.installRunning.Store(false)

	s.mu.Lock()
	last := s.last
	s.mu.Unlock()

	var info domain.UpdateInfo
	if last != nil && last.Available {
		info = *last
	} else {
		var err error
		if info, err = s.CheckForUpdate(ctx); err != nil {
			return err
		}
	}

	if !info.Available {
		return domain.ErrNoUpdate
	}
	if info.Asset == nil {
		return domain.ErrManualUpdate
	}
	if progress == nil {
		progress = func(int64, int64) {}
	}

	return s.installer.Install(ctx, *info.Asset, progress)
}

// ReleasePageURL returns the page of the latest known release, or the
// generic releases page if no check has completed yet.
func (s *UpdateService) ReleasePageURL() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.last != nil && s.last.ReleaseURL != "" {
		return s.last.ReleaseURL
	}
	return s.releasesURL
}
