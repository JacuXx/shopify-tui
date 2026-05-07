package app

import (
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/store"
	"github.com/JacuXx/shopify-cli/internal/views/addstore"
	"github.com/JacuXx/shopify-cli/internal/views/logs"
	"github.com/JacuXx/shopify-cli/internal/views/menu"
	"github.com/JacuXx/shopify-cli/internal/views/popup"
	"github.com/JacuXx/shopify-cli/internal/views/selectstore"
	"github.com/JacuXx/shopify-cli/internal/views/servers"
)

// Model es el estado central slim de la aplicación (patrón Elm Architecture).
type Model struct {
	vista domain.Vista
	ancho int
	alto  int

	// Sub-states por vista
	menu        menu.State
	addStore    addstore.State
	selectStore selectstore.State
	logsView    logs.State
	serversView servers.State
	popup       popup.State

	// Estado compartido entre vistas
	tiendas       []domain.Tienda
	tiendaParaDev domain.Tienda
	mensaje       string

	hayActualizacion bool
	versionNueva     string

	// Dependencias inyectadas
	storeRepo store.Repository
	serverMgr server.Manager
}
