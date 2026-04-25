package main

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
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
		
		// Ajustamos el tamaño de la lista considerando el banner (aprox 8-10 líneas)
		m.lista.SetSize(msg.Width-4, msg.Height-14)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":

			m.serverMgr.DetenerTodos()
			return m, tea.Quit
		case "q", "esc":

			switch m.vista {
			case VistaMenu:

				return m, nil
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

				items := crearListaMetodos()
				m.lista = crearLista(items, icons.Icons.Download+" Método de descarga", m.ancho, m.alto)
			case VistaSeleccionarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			case VistaSeleccionarModo:
				m.vista = VistaSeleccionarTienda
				m.mensaje = ""
				items := crearListaTiendas(m.tiendas, m.serverMgr)
				m.lista = crearLista(items, icons.Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
			case VistaLogs:

				m.vista = VistaSeleccionarModo
				tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)
				items := crearListaModos(m.tiendaParaDev, tieneServidor)
				titulo := icons.Icons.Server + " " + m.tiendaParaDev.Nombre
				if tieneServidor {
					titulo = icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
				}
				m.lista = crearLista(items, titulo, m.ancho, m.alto)
				m.mensaje = ""
			case VistaServidores:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			}
			return m, nil
		}

	case commands.ComandoTerminadoMsg:
		m.mensaje = msg.Resultado

		if msg.Tienda != nil {
			m.tiendas = append(m.tiendas, *msg.Tienda)
			m.storeRepo.GuardarTiendas(m.tiendas)
			m.mensaje = icons.IconSuccess("Tienda '" + msg.Tienda.Nombre + "' agregada correctamente")
			m.vista = domain.VistaMenu
			m.recrearMenuPrincipal()
			return m, nil
		}

		if msg.VolverAOpciones && m.tiendaParaDev.Nombre != "" {
			m.vista = domain.VistaLogs
			m.logsScroll = 0
			return m, tickCmd()
		}

		m.vista = domain.VistaMenu
		m.recrearMenuPrincipal()
		return m, nil

	case commands.ErrorMsg:
		m.mensaje = icons.IconError("Error: " + msg.Err.Error())
		return m, nil

	case tickMsg:

		if m.vista == VistaLogs {
			return m, tickCmd()
		}
		return m, nil
	}

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
	case VistaPopup:
		return m.updatePopup(msg)
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "a":
			return m, commands.EjecutarShopifyLogin()

		case "t":
			m.vista = VistaAgregarTienda
			m.inputNombre.SetValue("")
			m.inputURL.SetValue("")
			m.inputNombre.Focus()
			m.cursorInput = 0
			m.mensaje = ""
			m.tiendaTemporal = Tienda{}
			return m, nil

		case "d":
			if len(m.tiendas) == 0 {
				m.mensaje = icons.IconWarning("No hay tiendas. Agrega una primero.")
				return m, nil
			}
			m.vista = VistaSeleccionarTienda
			items := crearListaTiendas(m.tiendas, m.serverMgr)
			m.lista = crearLista(items, icons.Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
			m.mensaje = ""
			return m, nil

		case "v":
			m.vista = VistaServidores
			m.mensaje = ""
			return m, nil

		case "enter", "l":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			titulo := item.titulo

			switch {
			case strings.Contains(titulo, "Iniciar sesión"):
				return m, commands.EjecutarShopifyLogin()

			case strings.Contains(titulo, "Agregar tienda"):
				m.vista = VistaAgregarTienda
				m.inputNombre.SetValue("")
				m.inputURL.SetValue("")
				m.inputNombre.Focus()
				m.cursorInput = 0
				m.mensaje = ""
				m.tiendaTemporal = Tienda{}
				return m, nil

			case strings.Contains(titulo, "Desarrollo local"):
				if len(m.tiendas) == 0 {
					m.mensaje = icons.IconWarning("No hay tiendas guardadas. Agrega una primero.")
					return m, nil
				}
				m.vista = VistaSeleccionarTienda
				items := crearListaTiendas(m.tiendas, m.serverMgr)
				m.lista = crearLista(items, icons.Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
				m.mensaje = ""
				return m, nil

			case strings.Contains(titulo, "Servidores activos"):
				m.vista = VistaServidores
				m.mensaje = ""
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

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
				m.mensaje = icons.IconWarning("Por favor completa ambos campos")
				return m, nil
			}

			m.tiendaTemporal = Tienda{
				Nombre: nombre,
				URL:    url + ".myshopify.com",
			}

			m.vista = VistaSeleccionarMetodo
			items := crearListaMetodos()
			m.lista = crearLista(items, icons.Icons.Download+" Método de descarga", m.ancho, m.alto)
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

func (m Model) updateSeleccionarMetodo(msg tea.Msg) (tea.Model, tea.Cmd) {

	usarShopifyPull := func() (tea.Model, tea.Cmd) {
		directorio, err := m.storeRepo.CrearDirectorioTienda(m.tiendaTemporal.Nombre)
		if err != nil {
			m.mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return m, nil
		}
		m.tiendaTemporal.Metodo = domain.MetodoShopifyPull
		m.tiendaTemporal.Ruta = directorio
		return m, commands.EjecutarDescargaConExec(m.tiendaTemporal, directorio)
	}

	usarGitClone := func() (tea.Model, tea.Cmd) {
		directorio, err := m.storeRepo.CrearDirectorioTienda(m.tiendaTemporal.Nombre)
		if err != nil {
			m.mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return m, nil
		}
		m.tiendaTemporal.Metodo = domain.MetodoGitClone
		m.tiendaTemporal.Ruta = directorio
		m.vista = VistaInputGit
		m.inputGit.SetValue("")
		m.inputGit.Focus()
		m.mensaje = ""
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "s":
			return usarShopifyPull()
		case "g":
			return usarGitClone()

		case "enter":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			titulo := item.titulo
			if strings.Contains(titulo, "Shopify Pull") {
				return usarShopifyPull()
			} else if strings.Contains(titulo, "Git Clone") {
				return usarGitClone()
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m Model) updateInputGit(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			gitURL := strings.TrimSpace(m.inputGit.Value())
			if gitURL == "" {
				m.mensaje = icons.IconWarning("Por favor ingresa la URL del repositorio")
				return m, nil
			}
			if !m.gitURLConfirmada {
				// Primer Enter: mostrar confirmación
				m.gitURLConfirmada = true
				m.mensaje = icons.IconInfo("URL: " + gitURL + " — presiona Enter para confirmar")
				return m, nil
			}
			// Segundo Enter: clonar
			m.gitURLConfirmada = false
			m.tiendaTemporal.GitURL = gitURL
			return m, commands.EjecutarDescargaConExec(m.tiendaTemporal, m.tiendaTemporal.Ruta)
		case "esc", "q":
			m.gitURLConfirmada = false
		}
	}

	// Cualquier otra tecla resetea la confirmación pendiente
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() != "enter" {
			m.gitURLConfirmada = false
		}
	}

	var cmd tea.Cmd
	m.inputGit, cmd = m.inputGit.Update(msg)
	// Limpiar newlines que puedan quedar del paste
	if val := m.inputGit.Value(); strings.ContainsRune(val, '\n') || strings.ContainsRune(val, '\r') {
		m.inputGit.SetValue(strings.TrimSpace(val))
	}
	return m, cmd
}

func (m Model) updateSeleccionarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.lista.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.lista, cmd = m.lista.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		var indiceSeleccionado int = -1

		if len(key) == 1 && key >= "1" && key <= "9" {
			indiceSeleccionado = int(key[0] - '1')
			if indiceSeleccionado >= len(m.lista.VisibleItems()) {
				indiceSeleccionado = -1
			} else {
				m.lista.Select(indiceSeleccionado)
				indiceSeleccionado = m.lista.Index()
			}
		} else if key == "enter" || key == "l" {
			indiceSeleccionado = m.lista.Index()
		} else if key == "d" {
			if item, ok := m.lista.SelectedItem().(itemTienda); ok {
				tienda := item.tienda
				indiceReal := -1
				for i, t := range m.tiendas {
					if t.Nombre == tienda.Nombre && t.URL == tienda.URL {
						indiceReal = i
						break
					}
				}

				if indiceReal != -1 {
					nombreEliminada := tienda.Nombre
					rutaEliminada := tienda.Ruta

					if rutaEliminada != "" {
						if err := os.RemoveAll(rutaEliminada); err != nil {
							m.mensaje = icons.IconError("Error al eliminar archivos: " + err.Error())
							return m, nil
						}
					}

					m.tiendas = m.storeRepo.EliminarTienda(m.tiendas, indiceReal)

					if err := m.storeRepo.GuardarTiendas(m.tiendas); err != nil {
						m.mensaje = icons.IconError("Error al eliminar del registro: " + err.Error())
					} else {
						m.mensaje = icons.Icons.Delete + " Tienda '" + nombreEliminada + "' eliminada por completo"
					}

					if len(m.tiendas) == 0 {
						m.recrearMenuPrincipal()
						m.vista = domain.VistaMenu
						return m, nil
					}

					items := crearListaTiendas(m.tiendas, m.serverMgr)
					m.lista.SetItems(items)
				}
			}
			return m, nil
		}

		if indiceSeleccionado >= 0 {
			// Usar el item seleccionado de la lista (maneja filtrado automáticamente)
			if item, ok := m.lista.SelectedItem().(itemTienda); ok {
				tienda := item.tienda
				
				if !m.storeRepo.ExisteDirectorio(tienda.Ruta) {
					m.mensaje = icons.IconError("El directorio no existe: " + tienda.Ruta)
					return m, nil
				}

				m.tiendaParaDev = tienda
				tieneServidor := m.serverMgr.TieneServidorActivo(tienda.Nombre)

				if tieneServidor {
					m.vista = domain.VistaLogs
					m.logsScroll = 0
					m.mensaje = ""
					return m, tickCmd()
				}

				servidor, err := m.serverMgr.IniciarServidor(tienda)
				if err != nil {
					m.mensaje = icons.IconError(err.Error())
					m.vista = domain.VistaLogs
					m.logsScroll = 0
					return m, tickCmd()
				}

				m.mensaje = icons.IconSuccess("Servidor iniciado en " + servidor.URL)
				m.vista = domain.VistaLogs
				m.logsScroll = 0
				return m, tickCmd()
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m Model) updateSeleccionarModo(msg tea.Msg) (tea.Model, tea.Cmd) {
	tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)

	iniciarServidor := func() (tea.Model, tea.Cmd) {
		servidor, err := m.serverMgr.IniciarServidor(m.tiendaParaDev)
		if err != nil {
			m.mensaje = icons.IconError(err.Error())
			return m, nil
		}
		m.mensaje = icons.IconSuccess("Servidor iniciado en " + servidor.URL)
		m.vista = domain.VistaLogs
		m.logsScroll = 0
		return m, tickCmd()
	}

	verLogs := func() (tea.Model, tea.Cmd) {
		m.vista = domain.VistaLogs
		m.logsScroll = 0
		m.mensaje = ""
		return m, tickCmd()
	}

	detenerServidor := func() (tea.Model, tea.Cmd) {
		if err := m.serverMgr.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
			m.mensaje = icons.IconError(err.Error())
		} else {
			m.mensaje = icons.Icons.Stop + " Servidor detenido"
		}
		m.vista = domain.VistaMenu
		m.recrearMenuPrincipal()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "i":
			if !tieneServidor {
				return iniciarServidor()
			}
		case "l":
			if tieneServidor {
				return verLogs()
			}
		case "s":
			if tieneServidor {
				return detenerServidor()
			}
		case "p":
			return m, commands.EjecutarThemePull(m.tiendaParaDev)
		case "u":
			return m, commands.EjecutarThemePush(m.tiendaParaDev)
		case "e":
			return m, commands.EjecutarAbrirEditor(m.tiendaParaDev)
		case "t":
			return m, commands.EjecutarAbrirTerminal(m.tiendaParaDev)

		case "enter":
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			titulo := item.titulo

			switch {
			case strings.Contains(titulo, "Iniciar"):
				return iniciarServidor()

			case strings.Contains(titulo, "Ver logs"):
				return verLogs()

			case strings.Contains(titulo, "Detener"):
				return detenerServidor()

			case strings.Contains(titulo, "Pull"):
				return m, commands.EjecutarThemePull(m.tiendaParaDev)

			case strings.Contains(titulo, "Push"):
				return m, commands.EjecutarThemePush(m.tiendaParaDev)

			case strings.Contains(titulo, "Editor"):
				return m, commands.EjecutarAbrirEditor(m.tiendaParaDev)

			case strings.Contains(titulo, "Terminal"):
				return m, commands.EjecutarAbrirTerminal(m.tiendaParaDev)
			}
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m Model) updateServidores(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":

			servidores := m.serverMgr.ObtenerServidoresActivos()
			if len(servidores) == 0 {
				return m, nil
			}

			indice := m.lista.Index()
			if indice < 0 || indice >= len(servidores) {
				indice = 0
			}

			servidor := servidores[indice]
			if err := m.serverMgr.DetenerServidor(servidor.Tienda.Nombre); err != nil {
				m.mensaje = icons.IconError(err.Error())
			} else {
				m.mensaje = icons.IconSuccess("Servidor de '" + servidor.Tienda.Nombre + "' detenido")
			}
			return m, nil

		case "S":

			m.serverMgr.DetenerTodos()
			m.mensaje = "✅ Todos los servidores detenidos"
			return m, nil

		case "j", "down":

			servidores := m.serverMgr.ObtenerServidoresActivos()
			if len(servidores) > 0 {
				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil

		case "k", "up":

			servidores := m.serverMgr.ObtenerServidoresActivos()
			if len(servidores) > 0 {
				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) updateLogs(msg tea.Msg) (tea.Model, tea.Cmd) {
	servidor := m.serverMgr.ObtenerServidor(m.tiendaParaDev.Nombre)

	getMaxScroll := func() int {
		if servidor == nil {
			return 0
		}
		logs := servidor.ObtenerLogs()
		maxScroll := len(logs) - (m.alto - 8)
		if maxScroll < 0 {
			maxScroll = 0
		}
		return maxScroll
	}

	switch msg := msg.(type) {
	case tea.MouseMsg:

		if m.modoSeleccion {
			return m, nil
		}

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

		if m.modoSeleccion {
			if key == "v" {
				m.modoSeleccion = false
				m.mensaje = ""
				return m, tea.EnableMouseCellMotion
			}

			return m, nil
		}

		switch key {
		case " ", "m", "ctrl+p":
			m.vistaAnterior = VistaLogs
			m.vista = VistaPopup
			m.popupIndex = 0
			return m, nil

		case "v":

			m.modoSeleccion = true
			m.mensaje = icons.IconInfo("Modo selección ON - Usa Ctrl+Shift+C para copiar, 'v' para salir")

			return m, tea.DisableMouse

		case "ctrl+q":

			m.modoSeleccion = false
			m.vista = domain.VistaSeleccionarModo
			tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)
			items := crearListaModos(m.tiendaParaDev, tieneServidor)
			titulo := icons.Icons.Server + " " + m.tiendaParaDev.Nombre
			if tieneServidor {
				titulo = icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
			}
			m.lista = crearLista(items, titulo, m.ancho, m.alto)
			m.mensaje = ""

			return m, tea.EnableMouseCellMotion

		case "ctrl+s":

			if err := m.serverMgr.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = icons.IconError(err.Error())
			} else {
				m.mensaje = icons.Icons.Stop + " Servidor detenido"
			}
			return m, tickCmd()

		case "ctrl+g", "G":

			m.logsScroll = getMaxScroll()
			return m, nil

		case "ctrl+t", "g":

			m.logsScroll = 0
			return m, nil

		case "j", "down":

			maxScroll := getMaxScroll()
			m.logsScroll++
			if m.logsScroll > maxScroll {
				m.logsScroll = maxScroll
			}
			return m, nil

		case "k", "up":

			m.logsScroll--
			if m.logsScroll < 0 {
				m.logsScroll = 0
			}
			return m, nil

		case "pgdown", "ctrl+d":

			maxScroll := getMaxScroll()
			m.logsScroll += 10
			if m.logsScroll > maxScroll {
				m.logsScroll = maxScroll
			}
			return m, nil

		case "pgup", "ctrl+u":

			m.logsScroll -= 10
			if m.logsScroll < 0 {
				m.logsScroll = 0
			}
			return m, nil

		default:

			if servidor != nil && servidor.Activo {

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

					if len(key) == 1 {
						input = key
					}
				}

				if input != "" {
					if err := servidor.EnviarInput(input); err != nil {
						m.mensaje = icons.IconWarning("Error enviando input")
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func crearOpcionesPopup(tieneServidor bool) []itemMenu {
	opciones := []itemMenu{
		{titulo: icons.Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
		{titulo: icons.Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
		{titulo: icons.Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
		{titulo: icons.Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
	}

	if tieneServidor {
		opciones = append([]itemMenu{
			{titulo: icons.Icons.Stop + " Detener", desc: "Parar servidor", atajo: "s"},
		}, opciones...)
	}

	return opciones
}

func (m Model) updatePopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	tieneServidor := m.serverMgr.TieneServidorActivo(m.tiendaParaDev.Nombre)
	opciones := crearOpcionesPopup(tieneServidor)

	ejecutarOpcion := func(indice int) (tea.Model, tea.Cmd) {
		if indice < 0 || indice >= len(opciones) {
			return m, nil
		}

		opcion := opciones[indice]
		m.vista = domain.VistaLogs

		switch {
		case strings.Contains(opcion.titulo, "Detener"):
			if err := m.serverMgr.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = icons.IconError(err.Error())
			} else {
				m.mensaje = icons.Icons.Stop + " Servidor detenido"
				m.vista = domain.VistaMenu
				m.recrearMenuPrincipal()
			}
			return m, nil

		case strings.Contains(opcion.titulo, "Pull"):
			return m, commands.EjecutarThemePull(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Push"):
			return m, commands.EjecutarThemePush(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Editor"):
			return m, commands.EjecutarAbrirEditor(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Terminal"):
			return m, commands.EjecutarAbrirTerminal(m.tiendaParaDev)
		}

		return m, tickCmd()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "esc", "q", " ", "m", "ctrl+p":
			m.vista = VistaLogs
			return m, tickCmd()

		case "j", "down":
			m.popupIndex++
			if m.popupIndex >= len(opciones) {
				m.popupIndex = 0
			}
			return m, nil

		case "k", "up":
			m.popupIndex--
			if m.popupIndex < 0 {
				m.popupIndex = len(opciones) - 1
			}
			return m, nil

		case "enter", "l":
			return ejecutarOpcion(m.popupIndex)

		case "s":
			if tieneServidor {
				for i, op := range opciones {
					if strings.Contains(op.titulo, "Detener") {
						return ejecutarOpcion(i)
					}
				}
			}
		case "p":
			for i, op := range opciones {
				if strings.Contains(op.titulo, "Pull") {
					return ejecutarOpcion(i)
				}
			}
		case "u":
			for i, op := range opciones {
				if strings.Contains(op.titulo, "Push") {
					return ejecutarOpcion(i)
				}
			}
		case "e":
			for i, op := range opciones {
				if strings.Contains(op.titulo, "Editor") {
					return ejecutarOpcion(i)
				}
			}
		case "t":
			for i, op := range opciones {
				if strings.Contains(op.titulo, "Terminal") {
					return ejecutarOpcion(i)
				}
			}
		}
	}

	return m, nil
}
