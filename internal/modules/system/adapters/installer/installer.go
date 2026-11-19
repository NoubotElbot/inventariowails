// Package installer is the driven adapter that downloads a release asset
// and hands it to the operating system (installer, disk image or AppImage).
package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"inventario/internal/modules/system/core/domain"
	"inventario/internal/modules/system/core/ports"
)

var _ ports.Installer = (*Installer)(nil)

type Installer struct {
	http *http.Client
}

func New() *Installer {
	// No overall timeout: installers can be large and connections slow.
	// Cancellation is driven by the context.
	return &Installer{http: &http.Client{Timeout: 0}}
}

// Install downloads the asset to a temporary directory and launches it.
func (i *Installer) Install(ctx context.Context, asset domain.Asset, progress ports.ProgressFunc) error {
	if !supported() {
		return domain.ErrManualUpdate
	}

	path, err := i.download(ctx, asset, progress)
	if err != nil {
		return err
	}
	return launch(path)
}

func (i *Installer) download(ctx context.Context, asset domain.Asset, progress ports.ProgressFunc) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.DownloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "inventariowails-updater")

	resp, err := i.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("descargando actualización: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("descargando actualización: respuesta HTTP %d", resp.StatusCode)
	}

	total := resp.ContentLength
	if total <= 0 && asset.Size > 0 {
		total = asset.Size
	}
	if total <= 0 {
		total = -1
	}

	dir, err := os.MkdirTemp("", "inventario-update-*")
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, filepath.Base(asset.Name))

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}

	pw := &progressWriter{total: total, report: progress}
	_, copyErr := io.Copy(io.MultiWriter(f, pw), resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("descargando actualización: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.RemoveAll(dir)
		return "", closeErr
	}
	if total > 0 && pw.done != total {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("descarga incompleta: %d de %d bytes", pw.done, total)
	}

	progress(pw.done, total)
	return path, nil
}

// progressWriter counts written bytes and reports at most every 100ms.
type progressWriter struct {
	done, total int64
	report      ports.ProgressFunc
	lastReport  time.Time
}

func (p *progressWriter) Write(b []byte) (int, error) {
	p.done += int64(len(b))
	if now := time.Now(); now.Sub(p.lastReport) >= 100*time.Millisecond {
		p.lastReport = now
		p.report(p.done, p.total)
	}
	return len(b), nil
}
