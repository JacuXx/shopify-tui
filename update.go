// update.go - Manejo de eventos y lógica
// Este archivo contiene la función Update() que procesa todas las teclas y mensajes

package main

import (
	"strings"
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
				m.lista.Title = Icons.Download + " Método de descarga"
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
				m.lista.Title = Icons.Server + " Selecciona una tienda"
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
					m.lista.Title = Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
				} else {
					m.lista.Title = Icons.Server + " " + m.tiendaParaDev.Nombre
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
			m.mensaje = IconSuccess("Tienda '" + msg.tienda.Nombre + "' agregada correctamente")
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
				m.lista.Title = Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = Icons.Server + " " + m.tiendaParaDev.Nombre
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
		m.mensaje = IconError("Error: " + msg.err.Error())
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

			titulo := item.titulo
			
			// Usar strings.Contains para que funcione con cualquier set de iconos
			switch {
			case strings.Contains(titulo, "Iniciar sesión"):
				return m, ejecutarShopifyLogin()

			case strings.Contains(titulo, "Agregar tienda"):
				m.vista = VistaAgregarTienda
				m.inputNombre.SetValue("")
				m.inputURL.SetValue("")
				m.inputNombre.Focus()
				m.cursorInput = 0
				m.mensaje = ""
				m.tiendaTemporal = Tienda{}
				return m, nil

			case strings.Contains(titulo, "Opciones de desarrollo"):
				if len(m.tiendas) == 0 {
					m.mensaje = IconWarning("No hay tiendas guardadas. Agrega una primero.")
					return m, nil
				}
				m.vista = VistaSeleccionarTienda
				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = Icons.Server + " Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
				m.mensaje = ""
				return m, nil

			case strings.Contains(titulo, "Ver servidores"):
				m.vista = VistaServidores
				m.mensaje = ""
				return m, nil

			case strings.Contains(titulo, "Salir"):
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
				m.mensaje = IconWarning("Por favor completa ambos campos")
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
			m.lista.Title = Icons.Download + " Método de descarga"
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
				m.mensaje = IconError("Error al crear directorio: " + err.Error())
				return m, nil
			}

			titulo := item.titulo
			if strings.Contains(titulo, "Shopify Pull") {
				// Método: Shopify Pull
				m.tiendaTemporal.Metodo = MetodoShopifyPull
				m.tiendaTemporal.Ruta = directorio
				return m, ejecutarDescargaConExec(m.tiendaTemporal, directorio)

			} else if strings.Contains(titulo, "Git Clone") {
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
				m.mensaje = IconWarning("Por favor ingresa la URL del repositorio")
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
				m.mensaje = IconError("El directorio no existe: " + item.tienda.Ruta)
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
				m.lista.Title = Icons.ServerOn + " " + item.tienda.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = Icons.Server + " " + item.tienda.Nombre
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
					m.mensaje = IconError("Error al eliminar: " + err.Error())
				} else {
					m.mensaje = Icons.Delete + " Tienda '" + nombreEliminada + "' eliminada"
				}

				if len(m.tiendas) == 0 {
					m.vista = VistaMenu
					m.recrearMenuPrincipal()
					return m, nil
				}

				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = Icons.Server + " Selecciona una tienda"
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
			titulo := item.titulo

			switch {
			case strings.Contains(titulo, "Iniciar servidor"):
				// Iniciar servidor y mostrar logs
				servidor, err := gestor.IniciarServidor(m.tiendaParaDev)
				if err != nil {
					m.mensaje = IconError(err.Error())
					return m, nil
				}
				m.mensaje = IconSuccess("Servidor iniciado en " + servidor.URL)
				m.vista = VistaLogs
				m.logsScroll = 0
				return m, tickCmd()

			case strings.Contains(titulo, "Ver logs"):
				// Ir a ver los logs del servidor activo
				m.vista = VistaLogs
				m.logsScroll = 0
				m.mensaje = ""
				return m, tickCmd()

			case strings.Contains(titulo, "Detener servidor"):
				// Detener el servidor de esta tienda
				if err := gestor.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
					m.mensaje = IconError(err.Error())
				} else {
					m.mensaje = Icons.Stop + " Servidor de '" + m.tiendaParaDev.Nombre + "' detenido"
				}
				m.vista = VistaMenu
				m.recrearMenuPrincipal()
				return m, nil

			case strings.Contains(titulo, "Bajar cambios"):
				// Ejecutar theme pull
				return m, ejecutarThemePull(m.tiendaParaDev)

			case strings.Contains(titulo, "Pushear cambios"):
				// Ejecutar theme push
				return m, ejecutarThemePush(m.tiendaParaDev)

			case strings.Contains(titulo, "Abrir editor"):
				// Abrir VS Code (o el editor configurado)
				return m, ejecutarAbrirEditor(m.tiendaParaDev)

			case strings.Contains(titulo, "Abrir terminal"):
				// Abrir terminal en el directorio del tema
				return m, ejecutarAbrirTerminal(m.tiendaParaDev)
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
				m.mensaje = IconError(err.Error())
			} else {
				m.mensaje = IconSuccess("Servidor de '" + servidor.Tienda.Nombre + "' detenido")
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
	
	// Helper para calcular maxScroll
	getMaxScroll := func() int {
		if servidor == nil {
			return 0
		}
		logs := servidor.ObtenerLogs()
		maxScroll := len(logs) - (m.alto - 8) // Altura disponible para logs
		if maxScroll < 0 {
			maxScroll = 0
		}
		return maxScroll
	}
	
	switch msg := msg.(type) {
	case tea.MouseMsg:
		// En modo selección, ignorar eventos de mouse para permitir selección nativa
		if m.modoSeleccion {
			return m, nil
		}
		// Soporte para scroll con mouse
		switch msg.Type {
		case tea.MouseWheelUp:
			m.logsScroll -= 3
			if m.logsScroll < 0 {
				m.logsScroll = 0
			}
			return m, nil
		case tea.MouseWheelDown:
			maxScroll := getMaxScroll()
			m.logsScroll += 3
			if m.logsScroll > maxScroll {
				m.logsScroll = maxScroll
			}
			return m, nil
		}
		
	case tea.KeyMsg:
		key := msg.String()
		
		// EN MODO SELECCIÓN: Solo responder a 'v' para salir, ignorar todo lo demás
		if m.modoSeleccion {
			if key == "v" {
				m.modoSeleccion = false
				m.mensaje = ""
				return m, tea.EnableMouseCellMotion
			}
			// Ignorar todas las demás teclas en modo selección
			return m, nil
		}
		
		// Teclas de control del TUI (solo cuando NO está en modo selección)
		switch key {
		case "v":
			// Toggle modo selección (permite copiar con mouse) - "v" como visual mode en vim
			m.modoSeleccion = true
			m.mensaje = IconInfo("Modo selección ON - Usa Ctrl+Shift+C para copiar, 'v' para salir")
			// Desactivar captura de mouse para permitir selección nativa
			return m, tea.DisableMouse
			
		case "ctrl+q":
			// Volver al menú de modos (el servidor sigue corriendo)
			m.modoSeleccion = false // Resetear modo selección
			m.vista = VistaSeleccionarModo
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
			items := crearListaModos(m.tiendaParaDev, tieneServidor)
			m.lista = list.New(items, list.NewDefaultDelegate(), 55, 10)
			if tieneServidor {
				m.lista.Title = Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
			} else {
				m.lista.Title = Icons.Server + " " + m.tiendaParaDev.Nombre
			}
			m.lista.SetShowStatusBar(false)
			m.lista.SetFilteringEnabled(false)
			m.mensaje = ""
			// Reactivar mouse al salir
			return m, tea.EnableMouseCellMotion

		case "ctrl+s":
			// Detener servidor desde la vista de logs
			if err := ObtenerGestor().DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = IconError(err.Error())
			} else {
				m.mensaje = Icons.Stop + " Servidor detenido"
			}
			return m, tickCmd()

		case "ctrl+g", "G":
			// Ir al final de los logs (G como en vim)
			m.logsScroll = getMaxScroll()
			return m, nil

		case "ctrl+t", "g":
			// Ir al inicio de los logs (g como en vim)
			m.logsScroll = 0
			return m, nil

		case "j", "down":
			// Scroll hacia abajo (una línea)
			maxScroll := getMaxScroll()
			m.logsScroll++
			if m.logsScroll > maxScroll {
				m.logsScroll = maxScroll
			}
			return m, nil

		case "k", "up":
			// Scroll hacia arriba (una línea)
			m.logsScroll--
			if m.logsScroll < 0 {
				m.logsScroll = 0
			}
			return m, nil

		case "pgdown", "ctrl+d":
			// Scroll rápido hacia abajo (media página)
			maxScroll := getMaxScroll()
			m.logsScroll += 10
			if m.logsScroll > maxScroll {
				m.logsScroll = maxScroll
			}
			return m, nil

		case "pgup", "ctrl+u":
			// Scroll rápido hacia arriba (media página)
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
						m.mensaje = IconWarning("Error enviando input")
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}
