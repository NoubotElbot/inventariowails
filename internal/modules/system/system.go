// Package system is the composition root for the system module: it wires
// the GitHub releases client and the platform installer (driven adapters)
// into the update service, which is exposed through the Wails handler.
package system

import (
	"fmt"

	"inventario/internal/modules/system/adapters/github"
	"inventario/internal/modules/system/adapters/handler"
	"inventario/internal/modules/system/adapters/installer"
	"inventario/internal/modules/system/core/services"
	"inventario/internal/platform/buildinfo"
)

const (
	repoOwner = "NoubotElbot"
	repoName  = "inventariowails"
)

// Module holds the public surface of the system module.
type Module struct {
	Handler *handler.SystemHandler
}

func New() *Module {
	svc := services.NewUpdateService(
		buildinfo.Version,
		fmt.Sprintf("https://github.com/%s/%s/releases/latest", repoOwner, repoName),
		github.NewReleaseClient(repoOwner, repoName),
		installer.New(),
	)

	return &Module{Handler: handler.NewSystemHandler(svc)}
}
