package services_test

import (
	"context"
	"errors"
	"testing"

	"inventario/internal/modules/system/core/domain"
	"inventario/internal/modules/system/core/ports"
	"inventario/internal/modules/system/core/services"
)

const releasesURL = "https://example.com/releases/latest"

var errNet = errors.New("network error")

type mockReleases struct {
	release domain.Release
	err     error
	calls   int
}

func (m *mockReleases) LatestRelease(context.Context) (domain.Release, error) {
	m.calls++
	return m.release, m.err
}

type mockInstaller struct {
	installed []domain.Asset
	err       error
}

func (m *mockInstaller) Install(_ context.Context, asset domain.Asset, progress ports.ProgressFunc) error {
	progress(1, 1)
	m.installed = append(m.installed, asset)
	return m.err
}

func release(tag string) domain.Release {
	return domain.Release{
		Tag:   tag,
		Notes: "notas",
		URL:   "https://example.com/releases/" + tag,
		Assets: []domain.Asset{
			{Name: "inventariowails-amd64-installer.exe", DownloadURL: "https://example.com/x.exe", SHA256: "abc123"},
		},
	}
}

func newService(current string, rel *mockReleases, inst *mockInstaller, goos string) *services.UpdateService {
	return services.NewUpdateService(current, releasesURL, rel, inst, services.WithPlatform(goos, "amd64"))
}

func TestCurrentVersion(t *testing.T) {
	if got := newService("v1.2.3", &mockReleases{}, &mockInstaller{}, "windows").CurrentVersion(); got != "1.2.3" {
		t.Fatalf("got %q", got)
	}
	if got := newService("dev", &mockReleases{}, &mockInstaller{}, "windows").CurrentVersion(); got != "dev" {
		t.Fatalf("got %q", got)
	}
}

func TestCheckForUpdate_DevBuildSkipsRequest(t *testing.T) {
	rel := &mockReleases{release: release("v9.0.0")}
	info, err := newService("dev", rel, &mockInstaller{}, "windows").CheckForUpdate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Available || rel.calls != 0 {
		t.Fatalf("dev build: available=%v calls=%d", info.Available, rel.calls)
	}
}

func TestCheckForUpdate(t *testing.T) {
	tests := []struct {
		name, current, latest, goos string
		wantAvailable, wantAsset    bool
	}{
		{"newer release", "v1.0.0", "v1.1.0", "windows", true, true},
		{"same version", "v1.1.0", "v1.1.0", "windows", false, false},
		{"older release", "v1.2.0", "v1.1.0", "windows", false, false},
		{"newer without installer", "v1.0.0", "v1.1.0", "freebsd", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := newService(tt.current, &mockReleases{release: release(tt.latest)}, &mockInstaller{}, tt.goos).
				CheckForUpdate(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if info.Available != tt.wantAvailable {
				t.Fatalf("available = %v, want %v", info.Available, tt.wantAvailable)
			}
			if (info.Asset != nil) != tt.wantAsset {
				t.Fatalf("asset = %v, want present=%v", info.Asset, tt.wantAsset)
			}
			if info.LatestVersion != tt.latest[1:] {
				t.Fatalf("latest = %q", info.LatestVersion)
			}
		})
	}
}

func TestCheckForUpdate_AssetWithoutChecksumIsManual(t *testing.T) {
	rel := release("v1.1.0")
	rel.Assets[0].SHA256 = ""

	info, err := newService("v1.0.0", &mockReleases{release: rel}, &mockInstaller{}, "windows").
		CheckForUpdate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !info.Available || info.Asset != nil {
		t.Fatalf("available=%v asset=%v, want available without asset", info.Available, info.Asset)
	}
}

func TestCheckForUpdate_ProviderError(t *testing.T) {
	_, err := newService("v1.0.0", &mockReleases{err: errNet}, &mockInstaller{}, "windows").
		CheckForUpdate(context.Background())
	if !errors.Is(err, errNet) {
		t.Fatalf("err = %v, want %v", err, errNet)
	}
}

func TestInstallUpdate(t *testing.T) {
	inst := &mockInstaller{}
	rel := &mockReleases{release: release("v1.1.0")}
	svc := newService("v1.0.0", rel, inst, "windows")

	if _, err := svc.CheckForUpdate(context.Background()); err != nil {
		t.Fatal(err)
	}

	var reported bool
	if err := svc.InstallUpdate(context.Background(), func(int64, int64) { reported = true }); err != nil {
		t.Fatal(err)
	}
	if len(inst.installed) != 1 || inst.installed[0].Name != "inventariowails-amd64-installer.exe" {
		t.Fatalf("installed = %v", inst.installed)
	}
	if !reported {
		t.Fatal("progress not reported")
	}
	if rel.calls != 1 {
		t.Fatalf("expected cached check, got %d calls", rel.calls)
	}
}

func TestInstallUpdate_ChecksWhenNoCache(t *testing.T) {
	inst := &mockInstaller{}
	rel := &mockReleases{release: release("v1.1.0")}

	if err := newService("v1.0.0", rel, inst, "windows").InstallUpdate(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if rel.calls != 1 || len(inst.installed) != 1 {
		t.Fatalf("calls=%d installed=%d", rel.calls, len(inst.installed))
	}
}

func TestInstallUpdate_Errors(t *testing.T) {
	tests := []struct {
		name, current, goos string
		want                error
	}{
		{"no update", "v1.1.0", "windows", domain.ErrNoUpdate},
		{"no installer for platform", "v1.0.0", "freebsd", domain.ErrManualUpdate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &mockInstaller{}
			err := newService(tt.current, &mockReleases{release: release("v1.1.0")}, inst, tt.goos).
				InstallUpdate(context.Background(), nil)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			if len(inst.installed) != 0 {
				t.Fatal("installer must not run")
			}
		})
	}
}

func TestReleasePageURL(t *testing.T) {
	svc := newService("v1.0.0", &mockReleases{release: release("v1.1.0")}, &mockInstaller{}, "windows")
	if got := svc.ReleasePageURL(); got != releasesURL {
		t.Fatalf("before check: %q", got)
	}
	if _, err := svc.CheckForUpdate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := svc.ReleasePageURL(); got != "https://example.com/releases/v1.1.0" {
		t.Fatalf("after check: %q", got)
	}
}
