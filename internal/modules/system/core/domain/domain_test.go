package domain_test

import (
	"testing"

	"inventario/internal/modules/system/core/domain"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in      string
		want    domain.Version
		wantErr bool
	}{
		{in: "v1.2.3", want: domain.Version{Major: 1, Minor: 2, Patch: 3}},
		{in: "1.2.3", want: domain.Version{Major: 1, Minor: 2, Patch: 3}},
		{in: "v2.0", want: domain.Version{Major: 2}},
		{in: "v1.4.0-rc.1", want: domain.Version{Major: 1, Minor: 4}},
		{in: "v1.4.0+build5", want: domain.Version{Major: 1, Minor: 4}},
		{in: "dev", wantErr: true},
		{in: "", wantErr: true},
		{in: "v1.2.3.4", wantErr: true},
		{in: "v1.-2.3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := domain.ParseVersion(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseVersion(%q) = %v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseVersion(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("ParseVersion(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	v := func(s string) domain.Version {
		t.Helper()
		out, err := domain.ParseVersion(s)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}

	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.10.0", "1.9.9", 1},
		{"2.0.0", "1.99.99", 1},
		{"0.9.0", "1.0.0", -1},
	}
	for _, tt := range tests {
		if got := v(tt.a).Compare(v(tt.b)); got != tt.want {
			t.Errorf("%s.Compare(%s) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSelectAsset(t *testing.T) {
	assets := []domain.Asset{
		{Name: "inventariowails-amd64-installer.exe"},
		{Name: "inventariowails-arm64-installer.exe"},
		{Name: "inventariowails.dmg"},
		{Name: "inventariowails-x86_64.AppImage"},
		{Name: "inventariowails-aarch64.AppImage"},
		{Name: "inventariowails_1.2.0_amd64.deb"},
	}

	tests := []struct {
		goos, goarch string
		want         string
	}{
		{"windows", "amd64", "inventariowails-amd64-installer.exe"},
		{"windows", "arm64", "inventariowails-arm64-installer.exe"},
		{"darwin", "arm64", "inventariowails.dmg"},
		{"darwin", "amd64", "inventariowails.dmg"},
		{"linux", "amd64", "inventariowails-x86_64.AppImage"},
		{"linux", "arm64", "inventariowails-aarch64.AppImage"},
		{"linux", "riscv64", ""},
		{"freebsd", "amd64", ""},
	}

	for _, tt := range tests {
		t.Run(tt.goos+"/"+tt.goarch, func(t *testing.T) {
			got := domain.SelectAsset(assets, tt.goos, tt.goarch)
			switch {
			case tt.want == "" && got != nil:
				t.Fatalf("got %q, want nil", got.Name)
			case tt.want != "" && got == nil:
				t.Fatalf("got nil, want %q", tt.want)
			case got != nil && got.Name != tt.want:
				t.Fatalf("got %q, want %q", got.Name, tt.want)
			}
		})
	}
}
