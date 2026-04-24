package presentation

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.ancho = msg.Width
		m.alto = msg.Height
		m.lista.SetSize(msg.Width-4, msg.Height-14)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			// Detener todos los servidores
			m.container.ServerManager.StopAll()
			return m, tea.Quit

		case "q", "esc":
			switch m.vista {
			case VistaMenu:
				return m, nil
			case VistaAgregarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recreateMainMenu()
			case VistaSeleccionarMetodo:
				m.vista = VistaAgregarTienda
				m.mensaje = ""
			case VistaInputGit:
				m.vista = VistaSeleccionarMetodo
				m.mensaje = ""
				items := createDownloadMethods()
				m.lista = createList(items, "📥 Método de descarga", m.ancho, m.alto)
			case VistaSeleccionarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recreateMainMenu()
			case VistaSeleccionarModo:
				m.vista = VistaSeleccionarTienda
				m.mensaje = ""
				items := createStoreList(m.tiendas)
				m.lista = createList(items, "🖥️ Selecciona una tienda", m.ancho, m.alto)
			case VistaLogs:
				m.vista = VistaSeleccionarModo
				m.mensaje = ""
			case VistaServidores:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recreateMainMenu()
			}
			return m, nil
		}

	case addStoreResultMsg:
		m.handleAddStoreResult(msg)
		return m, nil

	case tickMsg:
		// Actualizar estado periódicamente
		return m, tickCmd()
	}

	// Delegar a la lista
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}
