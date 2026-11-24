package domain

import (
	"strings"
	"time"
)

// Asset is a downloadable file attached to a release.
type Asset struct {
	Name        string
	DownloadURL string
	Size        int64
	// SHA256 is the lowercase hex digest of the file. Empty when the
	// provider does not publish one; such assets are never auto-installed.
	SHA256 string
}

// Release is a published application version.
type Release struct {
	Tag         string
	Name        string
	Notes       string
	URL         string
	PublishedAt time.Time
	Assets      []Asset
}

// UpdateInfo describes whether a newer release is available for the
// running application and, if so, which asset installs it on this platform.
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	Notes          string
	ReleaseURL     string
	PublishedAt    time.Time
	// Asset is the installer for the current platform. It is nil when the
	// release has no compatible installer and the update must be manual.
	Asset *Asset
}

// SelectAsset picks the installer matching the given platform among the
// release assets, following the names produced by the release workflow:
//
//	windows: <app>-<arch>-installer.exe
//	darwin:  <app>.dmg
//	linux:   <app>-<arch>.AppImage
func SelectAsset(assets []Asset, goos, goarch string) *Asset {
	for i := range assets {
		name := strings.ToLower(assets[i].Name)
		switch goos {
		case "windows":
			if strings.HasSuffix(name, "-installer.exe") && strings.Contains(name, goarch) {
				return &assets[i]
			}
		case "darwin":
			if strings.HasSuffix(name, ".dmg") {
				return &assets[i]
			}
		case "linux":
			if strings.HasSuffix(name, ".appimage") && containsAny(name, linuxArchAliases(goarch)) {
				return &assets[i]
			}
		}
	}
	return nil
}

func linuxArchAliases(goarch string) []string {
	switch goarch {
	case "amd64":
		return []string{"x86_64", "amd64"}
	case "arm64":
		return []string{"aarch64", "arm64"}
	default:
		return []string{goarch}
	}
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
