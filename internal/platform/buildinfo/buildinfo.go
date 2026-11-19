// Package buildinfo holds metadata injected at build time.
package buildinfo

// Version is the application version. It defaults to "dev" and is overridden
// for release builds with:
//
//	-ldflags "-X inventario/internal/platform/buildinfo.Version=v1.2.3"
var Version = "dev"
