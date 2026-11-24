package github_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"inventario/internal/modules/system/adapters/github"
	"inventario/internal/modules/system/core/domain"
)

func TestLatestRelease(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/releases/latest" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("missing User-Agent")
		}
		_, _ = w.Write([]byte(`{
			"tag_name": "v1.2.0",
			"name": "v1.2.0",
			"body": "cambios",
			"html_url": "https://github.com/owner/repo/releases/tag/v1.2.0",
			"published_at": "2026-11-19T12:00:00Z",
			"assets": [
				{"name": "app-amd64-installer.exe", "browser_download_url": "https://dl/app.exe", "size": 42,
				 "digest": "sha256:ABCDEF0123456789abcdef0123456789abcdef0123456789abcdef0123456789"},
				{"name": "app.dmg", "browser_download_url": "https://dl/app.dmg", "size": 7},
				{"name": "app.deb", "browser_download_url": "https://dl/app.deb", "size": 7, "digest": "md5:abc"}
			]
		}`))
	}))
	defer srv.Close()

	rel, err := github.NewReleaseClient("owner", "repo").WithBaseURL(srv.URL).LatestRelease(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rel.Tag != "v1.2.0" || rel.Notes != "cambios" || rel.URL == "" || rel.PublishedAt.IsZero() {
		t.Fatalf("unexpected release: %+v", rel)
	}
	if len(rel.Assets) != 3 || rel.Assets[0].DownloadURL != "https://dl/app.exe" || rel.Assets[0].Size != 42 {
		t.Fatalf("unexpected assets: %+v", rel.Assets)
	}
	if want := "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"; rel.Assets[0].SHA256 != want {
		t.Fatalf("sha256 = %q, want %q", rel.Assets[0].SHA256, want)
	}
	if rel.Assets[1].SHA256 != "" || rel.Assets[2].SHA256 != "" {
		t.Fatalf("missing/non-sha256 digests must be empty: %+v", rel.Assets[1:])
	}
}

func TestLatestRelease_Errors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   error
	}{
		{"no releases", http.StatusNotFound, domain.ErrNoReleases},
		{"rate limited", http.StatusForbidden, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			_, err := github.NewReleaseClient("owner", "repo").WithBaseURL(srv.URL).LatestRelease(context.Background())
			if err == nil {
				t.Fatal("expected error")
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}
