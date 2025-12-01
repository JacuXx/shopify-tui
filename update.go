// update.go - Manejo de eventos y lógica
// Este archivo contiene la función Update() que procesa todas las teclas y mensajes

package main

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg es un mensaje para actualizar la vista de logs periódicamente
type tickMsg time.Time

// tickCmd retorna un comando que envía un tickMsg cada 500ms
func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init se ejecuta una vez al inicio
func (m Model) Init() tea.Cmd {
	return nil
}

// Update procesa todos los eventos
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.ancho = msg.Width
		m.alto = msg.Height
		m.lista.SetSize(msg.Width-4, msg.Height-6)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.vista == VistaMenu {
				return m, tea.Quit
			}
		case "esc":
			// Volver según la vista actual
			switch m.vista {
			case VistaAgregarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			case VistaSeleccionarMetodo:
				m.vista = VistaAgregarTienda
				m.mensaje = ""
			case VistaInputGit:
				m.vista = VistaSeleccionarMetodo
				m.mensaje = ""
				// Recrear lista de métodos
				items := crearListaMetodos()
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 10)
				m.lista.Title = "📦 Método de descarga"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
			case VistaSeleccionarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			case VistaSeleccionarModo:
				m.vista = VistaSeleccionarTienda
				m.mensaje = ""
				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = "🚀 Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
			case VistaLogs:
				// Volver al menú de modos de esta tienda
				m.vista = VistaSeleccionarModo
				gestor := ObtenerGestor()
				tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
				items := crearListaModos(m.tiendaParaDev, tieneServidor)
				m.lista = list.New(items, list.NewDefaultDelegate(), 55, 10)
				if tieneServidor {
					m.lista.Title = "🟢 " + m.tiendaParaDev.Nombre + " (servidor activo)"
				} else {
					m.lista.Title = "⚡ " + m.tiendaParaDev.Nombre
				}
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
				m.mensaje = ""
			case VistaServidores:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			}
			return m, nil
		}

	case comandoTerminadoMsg:
		m.mensaje = msg.resultado
		// Si viene una tienda, agregarla a la lista
		if msg.tienda != nil {
			m.tiendas = append(m.tiendas, *msg.tienda)
			guardarTiendas(m.tiendas)
			m.mensaje = "✅ Tienda '" + msg.tienda.Nombre + "' agregada correctamente"
		}
		
		// Decidir a dónde volver
		if msg.volverAOpciones && m.tiendaParaDev.Nombre != "" {
			// Volver a opciones de desarrollo de la tienda actual
			m.vista = VistaSeleccionarModo
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
			items := crearListaModos(m.tiendaParaDev, tieneServidor)
			m.lista = list.New(items, list.NewDefaultDelegate(), 55, 14)
			if tieneServidor {
				m.lista.Title = "🟢 " + m.tiendaParaDev.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = "🛠️ " + m.tiendaParaDev.Nombre
			}
			m.lista.SetShowStatusBar(false)
			m.lista.SetFilteringEnabled(false)
		} else {
			// Volver al menú principal
			m.vista = VistaMenu
			m.recrearMenuPrincipal()
		}
		return m, nil

	case errorMsg:
		m.mensaje = "❌ Error: " + msg.err.Error()
		return m, nil

	case tickMsg:
		// Solo refrescar si estamos en la vista de logs
		if m.vista == VistaLogs {
			return m, tickCmd()
		}
		return m, nil
	}

	// Delegar a la vista actual
	switch m.vista {
	case VistaMenu:
		return m.updateMenu(msg)
	case VistaAgregarTienda:
		return m.updateAgregarTienda(msg)
	case VistaSeleccionarMetodo:
		return m.updateSeleccionarMetodo(msg)
	case VistaInputGit:
		return m.updateInputGit(msg)
	case VistaSeleccionarTienda:
		return m.updateSeleccionarTienda(msg)
	case VistaSeleccionarModo:
		return m.updateSeleccionarModo(msg)
	case VistaLogs:
		return m.updateLogs(msg)
	case VistaServidores:
		return m.updateServidores(msg)
	}

	return m, nil
}

// updateMenu maneja el menú principal
func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			switch item.titulo {
			case "🔐 Iniciar sesión en Shopify":
				return m, ejecutarShopifyLogin()

			case "➕ Agregar tienda":
				m.vista = VistaAgregarTienda
				m.inputNombre.SetValue("")
				m.inputURL.SetValue("")
				m.inputNombre.Focus()
				m.cursorInput = 0
				m.mensaje = ""
				m.tiendaTemporal = Tienda{}
				return m, nil

			case "🛠️ Opciones de desarrollo":
				if len(m.tiendas) == 0 {
					m.mensaje = "⚠️ No hay tiendas guardadas. Agrega una primero."
					return m, nil
				}
				m.vista = VistaSeleccionarTienda
				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = "🛠️ Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
				m.mensaje = ""
				return m, nil

			case "📺 Ver servidores activos":
				m.vista = VistaServidores
				m.mensaje = ""
				return m, nil

			case "❌ Salir":
				// Detener todos los servidores antes de salir
				ObtenerGestor().DetenerTodos()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

// updateAgregarTienda maneja el formulario de nombre y URL
func (m Model) updateAgregarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			if m.cursorInput == 0 {
				m.cursorInput = 1
				m.inputNombre.Blur()
				m.inputURL.Focus()
			} else {
				m.cursorInput = 0
				m.inputURL.Blur()
				m.inputNombre.Focus()
			}
			return m, nil

		case "shift+tab", "up":
			if m.cursorInput == 1 {
				m.cursorInput = 0
				m.inputURL.Blur()
				m.inputNombre.Focus()
			} else {
				m.cursorInput = 1
				m.inputNombre.Blur()
				m.inputURL.Focus()
			}
			return m, nil

		case "enter":
			nombre := m.inputNombre.Value()
			url := m.inputURL.Value()

			if nombre == "" || url == "" {
				m.mensaje = "⚠️ Por favor completa ambos campos"
				return m, nil
			}

			// Guardar datos temporales y pasar a seleccionar método
			m.tiendaTemporal = Tienda{
				Nombre: nombre,
				URL:    url,
			}

			// Ir a seleccionar método
			m.vista = VistaSeleccionarMetodo
			items := crearListaMetodos()
			m.lista = list.New(items, list.NewDefaultDelegate(), 50, 10)
			m.lista.Title = "📦 Método de descarga"
			m.lista.SetShowStatusBar(false)
			m.lista.SetFilteringEnabled(false)
			m.mensaje = ""
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.cursorInput == 0 {
		m.inputNombre, cmd = m.inputNombre.Update(msg)
	} else {
		m.inputURL, cmd = m.inputURL.Update(msg)
	}
	return m, cmd
}

// updateSeleccionarMetodo maneja la selección de método de descarga
func (m Model) updateSeleccionarMetodo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			// Crear el directorio para la tienda
			directorio, err := crearDirectorioTienda(m.tiendaTemporal.Nombre)
			if err != nil {
				m.mensaje = "❌ Error al crear directorio: " + err.Error()
				return m, nil
			}

			if item.titulo == "📥 Shopify Pull" {
				// Método: Shopify Pull
				m.tiendaTemporal.Metodo = MetodoShopifyPull
				m.tiendaTemporal.Ruta = directorio
				return m, ejecutarDescargaConExec(m.tiendaTemporal, directorio)

			} else if item.titulo == "🔗 Git Clone" {
				// Método: Git - ir a pedir URL
				m.tiendaTemporal.Metodo = MetodoGitClone
				m.tiendaTemporal.Ruta = directorio
				m.vista = VistaInputGit
				m.inputGit.SetValue("")
				m.inputGit.Focus()
				m.mensaje = ""
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

// updateInputGit maneja el input de URL de git
func (m Model) updateInputGit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			gitURL := m.inputGit.Value()
			if gitURL == "" {
				m.mensaje = "⚠️ Por favor ingresa la URL del repositorio"
				return m, nil
			}

			// Guardar URL de git y ejecutar clone
			m.tiendaTemporal.GitURL = gitURL
			return m, ejecutarDescargaConExec(m.tiendaTemporal, m.tiendaTemporal.Ruta)
		}
	}

	var cmd tea.Cmd
	m.inputGit, cmd = m.inputGit.Update(msg)
	return m, cmd
}

// updateSeleccionarTienda maneja la selección de tienda para theme dev
func (m Model) updateSeleccionarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.lista.SelectedItem().(itemTienda)
			if !ok {
				return m, nil
			}

			// Verificar que el directorio existe
			if !existeDirectorio(item.tienda.Ruta) {
				m.mensaje = "❌ El directorio no existe: " + item.tienda.Ruta
				return m, nil
			}

			// Guardar tienda seleccionada y pasar al menú de modos
			m.tiendaParaDev = item.tienda
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(item.tienda.Nombre)

			// Ir a seleccionar modo
			m.vista = VistaSeleccionarModo
			items := crearListaModos(item.tienda, tieneServidor)
			m.lista = list.New(items, list.NewDefaultDelegate(), 55, 10)
			if tieneServidor {
				m.lista.Title = "🟢 " + item.tienda.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = "⚡ " + item.tienda.Nombre
			}
			m.lista.SetShowStatusBar(false)
			m.lista.SetFilteringEnabled(false)
			m.mensaje = ""
			return m, nil

		case "d":
			// Eliminar tienda
			indice := m.lista.Index()
			if indice >= 0 && indice < len(m.tiendas) {
				nombreEliminada := m.tiendas[indice].Nombre
				m.tiendas = eliminarTienda(m.tiendas, indice)

				if err := guardarTiendas(m.tiendas); err != nil {
					m.mensaje = "❌ Error al eliminar: " + err.Error()
				} else {
					m.mensaje = "🗑️ Tienda '" + nombreEliminada + "' eliminada"
				}

				if len(m.tiendas) == 0 {
					m.vista = VistaMenu
					m.recrearMenuPrincipal()
					return m, nil
				}

				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = "🚀 Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

// updateSeleccionarModo maneja la selección de modo para el servidor
func (m Model) updateSeleccionarModo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			gestor := ObtenerGestor()

			switch item.titulo {
			case "🚀 Iniciar servidor":
				// Iniciar servidor y mostrar logs
				servidor, err := gestor.IniciarServidor(m.tiendaParaDev)
				if err != nil {
					m.mensaje = "❌ " + err.Error()
					return m, nil
				}
				m.mensaje = "✅ Servidor iniciado en " + servidor.URL
				m.vista = VistaLogs
				m.logsScroll = 0
				return m, tickCmd()

			case "📺 Ver logs en vivo":
				// Ir a ver los logs del servidor activo
				m.vista = VistaLogs
				m.logsScroll = 0
				m.mensaje = ""
				return m, tickCmd()

			case "🛑 Detener servidor":
				// Detener el servidor de esta tienda
				if err := gestor.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
					m.mensaje = "❌ " + err.Error()
				} else {
					m.mensaje = "🛑 Servidor de '" + m.tiendaParaDev.Nombre + "' detenido"
				}
				m.vista = VistaMenu
				m.recrearMenuPrincipal()
				return m, nil

			case "📥 Bajar cambios":
				// Ejecutar theme pull
				return m, ejecutarThemePull(m.tiendaParaDev)

			case "📤 Pushear cambios":
				// Ejecutar theme push
				return m, ejecutarThemePush(m.tiendaParaDev)

			case "📝 Abrir editor de código":
				// Abrir VS Code (o el editor configurado)
				return m, ejecutarAbrirEditor(m.tiendaParaDev)
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

// updateServidores maneja la vista de servidores activos
func (m Model) updateServidores(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			// Detener el servidor seleccionado
			servidores := ObtenerGestor().ObtenerServidoresActivos()
			if len(servidores) == 0 {
				return m, nil
			}

			indice := m.lista.Index()
			if indice < 0 || indice >= len(servidores) {
				indice = 0
			}

			servidor := servidores[indice]
			if err := ObtenerGestor().DetenerServidor(servidor.Tienda.Nombre); err != nil {
				m.mensaje = "❌ " + err.Error()
			} else {
				m.mensaje = "✅ Servidor de '" + servidor.Tienda.Nombre + "' detenido"
			}
			return m, nil

		case "S":
			// Detener todos los servidores
			ObtenerGestor().DetenerTodos()
			m.mensaje = "✅ Todos los servidores detenidos"
			return m, nil

		case "j", "down":
			// Navegar abajo en la lista de servidores
			servidores := ObtenerGestor().ObtenerServidoresActivos()
			if len(servidores) > 0 {
				// Usar la lista para manejar la navegación
				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil

		case "k", "up":
			// Navegar arriba
			servidores := ObtenerGestor().ObtenerServidoresActivos()
			if len(servidores) > 0 {
				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil
		}
	}

	return m, nil
}

// updateLogs maneja la vista de logs en tiempo real
func (m Model) updateLogs(msg tea.Msg) (tea.Model, tea.Cmd) {
	servidor := ObtenerGestor().ObtenerServidor(m.tiendaParaDev.Nombre)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		
		// Teclas de control del TUI
		switch key {
		case "ctrl+q":
			// Volver al menú de modos (el servidor sigue corriendo)
			m.vista = VistaSeleccionarModo
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
			items := crearListaModos(m.tiendaParaDev, tieneServidor)
			m.lista = list.New(items, list.NewDefaultDelegate(), 55, 10)
			if tieneServidor {
				m.lista.Title = "🟢 " + m.tiendaParaDev.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = "⚡ " + m.tiendaParaDev.Nombre
			}
			m.lista.SetShowStatusBar(false)
			m.lista.SetFilteringEnabled(false)
			m.mensaje = ""
			return m, nil

		case "ctrl+s":
			// Detener servidor desde la vista de logs
			if err := ObtenerGestor().DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = "❌ " + err.Error()
			} else {
				m.mensaje = "🛑 Servidor detenido"
			}
			return m, tickCmd()

		case "ctrl+g":
			// Ir al final de los logs
			if servidor != nil {
				logs := servidor.ObtenerLogs()
				maxScroll := len(logs) - 20
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.logsScroll = maxScroll
			}
			return m, nil

		case "ctrl+t":
			// Ir al inicio de los logs
			m.logsScroll = 0
			return m, nil

		case "pgdown":
			// Scroll rápido hacia abajo
			if servidor != nil {
				logs := servidor.ObtenerLogs()
				maxScroll := len(logs) - 20
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.logsScroll += 10
				if m.logsScroll > maxScroll {
					m.logsScroll = maxScroll
				}
			}
			return m, nil

		case "pgup":
			// Scroll rápido hacia arriba
			m.logsScroll -= 10
			if m.logsScroll < 0 {
				m.logsScroll = 0
			}
			return m, nil

		default:
			// Todas las demás teclas se envían al proceso de Shopify
			if servidor != nil && servidor.Activo {
				// Convertir tecla a lo que espera el proceso
				var input string
				switch key {
				case "enter":
					input = "\n"
				case "space":
					input = " "
				case "tab":
					input = "\t"
				case "backspace":
					input = "\b"
				default:
					// Si es una tecla simple (letra, número, etc.)
					if len(key) == 1 {
						input = key
					}
				}
				
				if input != "" {
					if err := servidor.EnviarInput(input); err != nil {
						m.mensaje = "⚠️ Error enviando input"
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}
