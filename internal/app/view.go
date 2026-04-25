package app

import (
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/views/addstore"
	"github.com/JacuXx/shopify-cli/internal/views/logs"
	"github.com/JacuXx/shopify-cli/internal/views/menu"
	"github.com/JacuXx/shopify-cli/internal/views/popup"
	"github.com/JacuXx/shopify-cli/internal/views/selectstore"
	"github.com/JacuXx/shopify-cli/internal/views/servers"
)

// View es el dispatcher puro de renderizado (patrón Elm Architecture).
func (m Model) View() string {
	banner := components.RenderBanner()
	var contenido string

	switch m.vista {
	case domain.VistaMenu:
		return menu.View(m.menu, m.tiendas, m.serverMgr, m.hayActualizacion, m.versionNueva, m.mensaje)

	case domain.VistaAgregarTienda:
		contenido = addstore.ViewAgregarTienda(m.addStore)

	case domain.VistaSeleccionarMetodo:
		contenido = addstore.ViewSeleccionarMetodo(m.addStore)

	case domain.VistaInputGit:
		contenido = addstore.ViewInputGit(m.addStore)

	case domain.VistaSeleccionarTienda:
		contenido = selectstore.ViewSeleccionarTienda(m.selectStore)

	case domain.VistaSeleccionarModo:
		contenido = selectstore.ViewSeleccionarModo(m.selectStore, m.tiendaParaDev, m.serverMgr)

	case domain.VistaLogs:
		servidor := m.serverMgr.ObtenerServidor(m.tiendaParaDev.Nombre)
		contenido = logs.View(m.logsView, servidor, m.tiendaParaDev, m.alto, m.mensaje)

	case domain.VistaServidores:
		contenido = servers.View(m.serversView, m.serverMgr)

	case domain.VistaPopup:
		contenido = popup.View(m.popup, m.tiendaParaDev, m.serverMgr, m.ancho)

	default:
		return menu.View(m.menu, m.tiendas, m.serverMgr, m.hayActualizacion, m.versionNueva, m.mensaje)
	}

	return banner + "\n" + contenido
}
