// Package github is the driven adapter that reads published releases from
// the GitHub REST API.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"inventario/internal/modules/system/core/domain"
	"inventario/internal/modules/system/core/ports"
)

const defaultBaseURL = "https://api.github.com"

var _ ports.ReleaseProvider = (*ReleaseClient)(nil)

type ReleaseClient struct {
	owner, repo string
	baseURL     string
	http        *http.Client
}

// NewReleaseClient reads releases of a public repository (no token needed).
func NewReleaseClient(owner, repo string) *ReleaseClient {
	return &ReleaseClient{
		owner:   owner,
		repo:    repo,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// WithBaseURL points the client to another API host (used in tests).
func (c *ReleaseClient) WithBaseURL(baseURL string) *ReleaseClient {
	c.baseURL = baseURL
	return c
}

type releaseDTO struct {
	TagName     string     `json:"tag_name"`
	Name        string     `json:"name"`
	Body        string     `json:"body"`
	HTMLURL     string     `json:"html_url"`
	PublishedAt time.Time  `json:"published_at"`
	Assets      []assetDTO `json:"assets"`
}

type assetDTO struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// LatestRelease returns the most recent non-draft, non-prerelease release.
func (c *ReleaseClient) LatestRelease(ctx context.Context) (domain.Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", c.baseURL, c.owner, c.repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Release{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "inventariowails-updater")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.Release{}, fmt.Errorf("consultando releases: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return domain.Release{}, domain.ErrNoReleases
	case resp.StatusCode != http.StatusOK:
		return domain.Release{}, fmt.Errorf("consultando releases: respuesta HTTP %d", resp.StatusCode)
	}

	var dto releaseDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return domain.Release{}, fmt.Errorf("leyendo release: %w", err)
	}

	return toDomain(dto), nil
}

func toDomain(dto releaseDTO) domain.Release {
	assets := make([]domain.Asset, 0, len(dto.Assets))
	for _, a := range dto.Assets {
		assets = append(assets, domain.Asset{
			Name:        a.Name,
			DownloadURL: a.BrowserDownloadURL,
			Size:        a.Size,
		})
	}

	return domain.Release{
		Tag:         dto.TagName,
		Name:        dto.Name,
		Notes:       dto.Body,
		URL:         dto.HTMLURL,
		PublishedAt: dto.PublishedAt,
		Assets:      assets,
	}
}
