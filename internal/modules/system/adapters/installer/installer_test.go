package installer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"inventario/internal/modules/system/core/domain"
)

func TestDownload(t *testing.T) {
	payload := strings.Repeat("x", 1<<16)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	var lastDone, lastTotal int64
	path, err := New().download(context.Background(), domain.Asset{
		Name:        "app-amd64-installer.exe",
		DownloadURL: srv.URL,
		Size:        int64(len(payload)),
	}, func(done, total int64) { lastDone, lastTotal = done, total })
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(filepath.Dir(path))

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != payload {
		t.Fatalf("downloaded %d bytes, want %d", len(got), len(payload))
	}
	if filepath.Base(path) != "app-amd64-installer.exe" {
		t.Fatalf("unexpected file name %q", path)
	}
	if lastDone != int64(len(payload)) || lastTotal != int64(len(payload)) {
		t.Fatalf("final progress = %d/%d", lastDone, lastTotal)
	}
}

func TestDownload_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := New().download(context.Background(), domain.Asset{Name: "a.exe", DownloadURL: srv.URL}, func(int64, int64) {})
	if err == nil {
		t.Fatal("expected error")
	}
}
