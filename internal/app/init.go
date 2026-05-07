package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/store"
	"github.com/JacuXx/shopify-cli/internal/version"
	"github.com/JacuXx/shopify-cli/internal/views/addstore"
	"github.com/JacuXx/shopify-cli/internal/views/logs"
	"github.com/JacuXx/shopify-cli/internal/views/menu"
	"github.com/JacuXx/shopify-cli/internal/views/popup"
	"github.com/JacuXx/shopify-cli/internal/views/selectstore"
	"github.com/JacuXx/shopify-cli/internal/views/servers"
)

// New crea el modelo inicial de la aplicación con todas las dependencias inyectadas.
func New() (tea.Model, tea.Cmd) {
	storeRepo := store.NewJSONRepository()
	serverMgr := server.NewProcessManager()
	tiendas, _ := storeRepo.CargarTiendas() //nolint:errcheck // app arranca con lista vacía si el archivo no existe
	hayUpdate, versionNew := version.VerificarActualizacion()

	m := Model{
		vista:            domain.VistaMenu,
		menu:             menu.NewState(0, 0),
		addStore:         addstore.NewState(),
		selectStore:      selectstore.NewState(),
		logsView:         logs.NewState(),
		serversView:      servers.NewState(),
		popup:            popup.NewState(),
		tiendas:          tiendas,
		hayActualizacion: hayUpdate,
		versionNueva:     versionNew,
		storeRepo:        storeRepo,
		serverMgr:        serverMgr,
	}
	return m, nil
}

// Init implementa tea.Model. No hay comandos iniciales.
func (m Model) Init() tea.Cmd {
	return nil
}
