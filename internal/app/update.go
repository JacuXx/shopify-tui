package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/views/addstore"
	"github.com/JacuXx/shopify-cli/internal/views/logs"
	"github.com/JacuXx/shopify-cli/internal/views/menu"
	"github.com/JacuXx/shopify-cli/internal/views/popup"
	"github.com/JacuXx/shopify-cli/internal/views/selectstore"
	"github.com/JacuXx/shopify-cli/internal/views/servers"
)

// Update es el dispatcher central de mensajes (patrón Elm Architecture).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.ancho = msg.Width
		m.alto = msg.Height
		m.menu.Rebuild(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			m.serverMgr.DetenerTodos()
			return m, tea.Quit

		case "q", "esc":
			switch m.vista {
			case domain.VistaMenu:
				return m, nil

			case domain.VistaAgregarTienda:
				m.vista = domain.VistaMenu
				m.mensaje = ""
				m.addStore = addstore.NewState()
				m.menu.Rebuild(m.ancho, m.alto)

			case domain.VistaSeleccionarMetodo:
				m.vista = domain.VistaAgregarTienda
				m.addStore.Mensaje = ""

			case domain.VistaInputGit:
				m.vista = domain.VistaSeleccionarMetodo
				m.addStore.GitURLConfirmada = false
				m.addStore.Mensaje = ""

			case domain.VistaSeleccionarTienda:
				m.vista = domain.VistaMenu
				m.mensaje = ""
				m.menu.Rebuild(m.ancho, m.alto)

			case domain.VistaSeleccionarModo:
				m.vista = domain.VistaSeleccionarTienda
				m.mensaje = ""
				items := components.CrearListaTiendas(m.tiendas, m.serverMgr)
				m.selectStore.Lista = components.CrearLista(items, icons.Icons.Server+" Selecciona una tienda", m.ancho, m.alto)

			case domain.VistaLogs:
				m.vista = domain.VistaSeleccionarModo
				tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)
				items := components.CrearListaModos(m.tiendaParaDev, tieneServidor)
				titulo := icons.Icons.Server + " " + m.tiendaParaDev.Nombre
				if tieneServidor {
					titulo = icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
				}
				m.selectStore.Lista = components.CrearLista(items, titulo, m.ancho, m.alto)
				m.mensaje = ""

			case domain.VistaServidores:
				m.vista = domain.VistaMenu
				m.mensaje = ""
				m.menu.Rebuild(m.ancho, m.alto)
			}
			return m, nil
		}

	case commands.ComandoTerminadoMsg:
		m.mensaje = msg.Resultado

		if msg.Tienda != nil {
			m.tiendas = append(m.tiendas, *msg.Tienda)
			if err := m.storeRepo.GuardarTiendas(m.tiendas); err != nil {
				m.mensaje = icons.IconError("Error al guardar tienda: " + err.Error())
				return m, nil
			}
			m.mensaje = icons.IconSuccess("Tienda '" + msg.Tienda.Nombre + "' agregada correctamente")
			m.vista = domain.VistaMenu
			m.menu.Rebuild(m.ancho, m.alto)
			return m, nil
		}

		if msg.VolverAOpciones && m.tiendaParaDev.Nombre != "" {
			m.vista = domain.VistaLogs
			m.logsView.Scroll = 0
			return m, logs.TickCmd()
		}

		m.vista = domain.VistaMenu
		m.menu.Rebuild(m.ancho, m.alto)
		return m, nil

	case commands.ErrorMsg:
		m.mensaje = icons.IconError("Error: " + msg.Err.Error())
		return m, nil

	case logs.TickMsg:
		if m.vista == domain.VistaLogs {
			return m, logs.TickCmd()
		}
		return m, nil

	case domain.NavMsg:
		return m.handleNav(msg)
	}

	// Delegar a sub-update según vista activa
	switch m.vista {
	case domain.VistaMenu:
		// Interceptar NavMsg de "d" sin tiendas para mostrar aviso en lugar de navegar
		newMenu, cmd := menu.Update(m.menu, msg, m.tiendas)
		m.menu = newMenu
		// Si el sub-update devuelve NavTo(VistaMenu) estando ya en menu, es señal de "no hay tiendas"
		return m, cmd

	case domain.VistaAgregarTienda:
		newAdd, cmd := addstore.UpdateAgregarTienda(m.addStore, msg)
		m.addStore = newAdd
		return m, cmd

	case domain.VistaSeleccionarMetodo:
		newAdd, cmd := addstore.UpdateSeleccionarMetodo(m.addStore, msg, m.storeRepo)
		m.addStore = newAdd
		return m, cmd

	case domain.VistaInputGit:
		newAdd, cmd := addstore.UpdateInputGit(m.addStore, msg)
		m.addStore = newAdd
		return m, cmd

	case domain.VistaSeleccionarTienda:
		newSel, cmd, tiendas := selectstore.UpdateSeleccionarTienda(m.selectStore, msg, m.tiendas, m.serverMgr, m.storeRepo, m.ancho, m.alto)
		m.selectStore = newSel
		m.tiendas = tiendas
		return m, cmd

	case domain.VistaSeleccionarModo:
		newSel, cmd := selectstore.UpdateSeleccionarModo(m.selectStore, msg, m.tiendaParaDev, m.serverMgr, m.ancho, m.alto)
		m.selectStore = newSel
		return m, cmd

	case domain.VistaLogs:
		servidor := m.serverMgr.ObtenerServidor(m.tiendaParaDev.Nombre)
		newLogs, cmd := logs.Update(m.logsView, msg, servidor, m.tiendaParaDev, m.serverMgr, m.alto)
		m.logsView = newLogs
		return m, cmd

	case domain.VistaServidores:
		newSrv, cmd := servers.Update(m.serversView, msg, m.serverMgr)
		m.serversView = newSrv
		return m, cmd

	case domain.VistaPopup:
		newPopup, cmd := popup.Update(m.popup, msg, m.tiendaParaDev, m.serverMgr)
		m.popup = newPopup
		return m, cmd
	}

	return m, nil
}

// handleNav procesa los mensajes de navegación emitidos por los sub-updates.
func (m Model) handleNav(msg domain.NavMsg) (tea.Model, tea.Cmd) {
	target := msg.Target

	switch target {
	case domain.VistaMenu:
		// Caso especial: menu/update emite NavTo(VistaMenu) cuando "d" y no hay tiendas
		if m.vista == domain.VistaMenu {
			m.mensaje = icons.IconWarning("No hay tiendas. Agrega una primero.")
			return m, nil
		}
		m.vista = domain.VistaMenu
		m.menu.Rebuild(m.ancho, m.alto)
		return m, nil

	case domain.VistaAgregarTienda:
		m.vista = domain.VistaAgregarTienda
		m.addStore = addstore.NewState()
		m.addStore.InputNombre.Focus()
		m.mensaje = ""
		return m, nil

	case domain.VistaSeleccionarMetodo:
		m.vista = domain.VistaSeleccionarMetodo
		m.addStore.Mensaje = ""
		return m, nil

	case domain.VistaInputGit:
		m.vista = domain.VistaInputGit
		return m, nil

	case domain.VistaSeleccionarTienda:
		m.vista = domain.VistaSeleccionarTienda
		items := components.CrearListaTiendas(m.tiendas, m.serverMgr)
		m.selectStore.Lista = components.CrearLista(items, icons.Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
		m.selectStore.Mensaje = ""
		m.mensaje = ""
		return m, nil

	case domain.VistaSeleccionarModo:
		m.vista = domain.VistaSeleccionarModo
		tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)
		items := components.CrearListaModos(m.tiendaParaDev, tieneServidor)
		titulo := icons.Icons.Server + " " + m.tiendaParaDev.Nombre
		if tieneServidor {
			titulo = icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
		}
		m.selectStore.Lista = components.CrearLista(items, titulo, m.ancho, m.alto)
		return m, nil

	case domain.VistaLogs:
		// Puede traer una tienda en el contexto (desde VistaSeleccionarTienda)
		if msg.Tienda != nil {
			m.tiendaParaDev = *msg.Tienda
		}
		m.vista = domain.VistaLogs
		m.logsView.Scroll = 0
		if m.selectStore.Mensaje != "" {
			m.mensaje = m.selectStore.Mensaje
		}
		return m, logs.TickCmd()

	case domain.VistaServidores:
		m.vista = domain.VistaServidores
		m.mensaje = ""
		return m, nil

	case domain.VistaPopup:
		m.vista = domain.VistaPopup
		m.popup.Index = 0
		return m, nil
	}

	return m, nil
}
