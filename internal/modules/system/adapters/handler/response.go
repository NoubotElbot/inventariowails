package handler

import (
	"inventario/internal/modules/system/core/domain"
	"inventario/internal/platform/responses"
)

type UpdateInfoResponse struct {
	Available      bool   `json:"available"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Notes          string `json:"notes"`
	ReleaseURL     string `json:"release_url"`
	PublishedAt    string `json:"published_at"`
	// AutoInstall is false when the platform has no compatible installer and
	// the user must download the update manually.
	AutoInstall bool `json:"auto_install"`
}

// UpdateProgress is the payload of UpdateProgressEvent. Total is -1 when
// the size is unknown.
type UpdateProgress struct {
	Downloaded int64 `json:"downloaded"`
	Total      int64 `json:"total"`
}

func toUpdateInfoResponse(info domain.UpdateInfo) UpdateInfoResponse {
	res := UpdateInfoResponse{
		Available:      info.Available,
		CurrentVersion: info.CurrentVersion,
		LatestVersion:  info.LatestVersion,
		Notes:          info.Notes,
		ReleaseURL:     info.ReleaseURL,
		AutoInstall:    info.Asset != nil,
	}
	if !info.PublishedAt.IsZero() {
		res.PublishedAt = responses.FormatTime(info.PublishedAt)
	}
	return res
}
