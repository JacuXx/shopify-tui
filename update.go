package main

import (
	"strings"
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
		m.lista.SetSize(msg.Width-4, msg.Height-6)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":

			ObtenerGestor().DetenerTodos()
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
				m.lista = crearLista(items, Icons.Download+" Método de descarga", m.ancho, m.alto)
			case VistaSeleccionarTienda:
				m.vista = VistaMenu
				m.mensaje = ""
				m.recrearMenuPrincipal()
			case VistaSeleccionarModo:
				m.vista = VistaSeleccionarTienda
				m.mensaje = ""
				items := crearListaTiendas(m.tiendas)
				m.lista = crearLista(items, Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
			case VistaLogs:

				m.vista = VistaSeleccionarModo
				gestor := ObtenerGestor()
				tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
				items := crearListaModos(m.tiendaParaDev, tieneServidor)
				titulo := Icons.Server + " " + m.tiendaParaDev.Nombre
				if tieneServidor {
					titulo = Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
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

	case comandoTerminadoMsg:
		m.mensaje = msg.resultado

		if msg.tienda != nil {
			m.tiendas = append(m.tiendas, *msg.tienda)
			guardarTiendas(m.tiendas)
			m.mensaje = IconSuccess("Tienda '" + msg.tienda.Nombre + "' agregada correctamente")
			m.vista = VistaMenu
			m.recrearMenuPrincipal()
			return m, nil
		}

		if msg.volverAOpciones && m.tiendaParaDev.Nombre != "" {
			m.vista = VistaLogs
			m.logsScroll = 0
			return m, tickCmd()
		}

		m.vista = VistaMenu
		m.recrearMenuPrincipal()
		return m, nil

	case errorMsg:
		m.mensaje = IconError("Error: " + msg.err.Error())
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
			return m, ejecutarShopifyLogin()

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
				m.mensaje = IconWarning("No hay tiendas. Agrega una primero.")
				return m, nil
			}
			m.vista = VistaSeleccionarTienda
			items := crearListaTiendas(m.tiendas)
			m.lista = crearLista(items, Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
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

			case strings.Contains(titulo, "Desarrollo local"):
				if len(m.tiendas) == 0 {
					m.mensaje = IconWarning("No hay tiendas guardadas. Agrega una primero.")
					return m, nil
				}
				m.vista = VistaSeleccionarTienda
				items := crearListaTiendas(m.tiendas)
				m.lista = crearLista(items, Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
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
				m.mensaje = IconWarning("Por favor completa ambos campos")
				return m, nil
			}

			m.tiendaTemporal = Tienda{
				Nombre: nombre,
				URL:    url + ".myshopify.com",
			}

			m.vista = VistaSeleccionarMetodo
			items := crearListaMetodos()
			m.lista = crearLista(items, Icons.Download+" Método de descarga", m.ancho, m.alto)
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
		directorio, err := crearDirectorioTienda(m.tiendaTemporal.Nombre)
		if err != nil {
			m.mensaje = IconError("Error al crear directorio: " + err.Error())
			return m, nil
		}
		m.tiendaTemporal.Metodo = MetodoShopifyPull
		m.tiendaTemporal.Ruta = directorio
		return m, ejecutarDescargaConExec(m.tiendaTemporal, directorio)
	}

	usarGitClone := func() (tea.Model, tea.Cmd) {
		directorio, err := crearDirectorioTienda(m.tiendaTemporal.Nombre)
		if err != nil {
			m.mensaje = IconError("Error al crear directorio: " + err.Error())
			return m, nil
		}
		m.tiendaTemporal.Metodo = MetodoGitClone
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l":
			gitURL := m.inputGit.Value()
			if gitURL == "" {
				m.mensaje = IconWarning("Por favor ingresa la URL del repositorio")
				return m, nil
			}

			m.tiendaTemporal.GitURL = gitURL
			return m, ejecutarDescargaConExec(m.tiendaTemporal, m.tiendaTemporal.Ruta)
		}
	}

	var cmd tea.Cmd
	m.inputGit, cmd = m.inputGit.Update(msg)
	return m, cmd
}

func (m Model) updateSeleccionarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		var indiceSeleccionado int = -1

		if len(key) == 1 && key >= "1" && key <= "9" {
			indiceSeleccionado = int(key[0] - '1')
		} else if key == "enter" || key == "l" {
			indiceSeleccionado = m.lista.Index()
		} else if key == "d" {
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
				m.lista = crearLista(items, Icons.Server+" Selecciona una tienda", m.ancho, m.alto)
			}
			return m, nil
		}

		if indiceSeleccionado >= 0 && indiceSeleccionado < len(m.tiendas) {
			tienda := m.tiendas[indiceSeleccionado]

			if !existeDirectorio(tienda.Ruta) {
				m.mensaje = IconError("El directorio no existe: " + tienda.Ruta)
				return m, nil
			}

			m.tiendaParaDev = tienda
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(tienda.Nombre)

			if tieneServidor {
				m.vista = VistaLogs
				m.logsScroll = 0
				m.mensaje = ""
				return m, tickCmd()
			}

			servidor, err := gestor.IniciarServidor(tienda)
			if err != nil {
				m.mensaje = IconError(err.Error())
				m.vista = VistaLogs
				m.logsScroll = 0
				return m, tickCmd()
			}

			m.mensaje = IconSuccess("Servidor iniciado en " + servidor.URL)
			m.vista = VistaLogs
			m.logsScroll = 0
			return m, tickCmd()
		}
	}

	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

func (m Model) updateSeleccionarModo(msg tea.Msg) (tea.Model, tea.Cmd) {
	gestor := ObtenerGestor()
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	iniciarServidor := func() (tea.Model, tea.Cmd) {
		servidor, err := gestor.IniciarServidor(m.tiendaParaDev)
		if err != nil {
			m.mensaje = IconError(err.Error())
			return m, nil
		}
		m.mensaje = IconSuccess("Servidor iniciado en " + servidor.URL)
		m.vista = VistaLogs
		m.logsScroll = 0
		return m, tickCmd()
	}

	verLogs := func() (tea.Model, tea.Cmd) {
		m.vista = VistaLogs
		m.logsScroll = 0
		m.mensaje = ""
		return m, tickCmd()
	}

	detenerServidor := func() (tea.Model, tea.Cmd) {
		if err := gestor.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
			m.mensaje = IconError(err.Error())
		} else {
			m.mensaje = Icons.Stop + " Servidor detenido"
		}
		m.vista = VistaMenu
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
			return m, ejecutarThemePull(m.tiendaParaDev)
		case "u":
			return m, ejecutarThemePush(m.tiendaParaDev)
		case "e":
			return m, ejecutarAbrirEditor(m.tiendaParaDev)
		case "t":
			return m, ejecutarAbrirTerminal(m.tiendaParaDev)

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
				return m, ejecutarThemePull(m.tiendaParaDev)

			case strings.Contains(titulo, "Push"):
				return m, ejecutarThemePush(m.tiendaParaDev)

			case strings.Contains(titulo, "Editor"):
				return m, ejecutarAbrirEditor(m.tiendaParaDev)

			case strings.Contains(titulo, "Terminal"):
				return m, ejecutarAbrirTerminal(m.tiendaParaDev)
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

			ObtenerGestor().DetenerTodos()
			m.mensaje = "✅ Todos los servidores detenidos"
			return m, nil

		case "j", "down":

			servidores := ObtenerGestor().ObtenerServidoresActivos()
			if len(servidores) > 0 {

				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil

		case "k", "up":

			servidores := ObtenerGestor().ObtenerServidoresActivos()
			if len(servidores) > 0 {
				m.lista, _ = m.lista.Update(msg)
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) updateLogs(msg tea.Msg) (tea.Model, tea.Cmd) {
	servidor := ObtenerGestor().ObtenerServidor(m.tiendaParaDev.Nombre)

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
			m.mensaje = IconInfo("Modo selección ON - Usa Ctrl+Shift+C para copiar, 'v' para salir")

			return m, tea.DisableMouse

		case "ctrl+q":

			m.modoSeleccion = false
			m.vista = VistaSeleccionarModo
			gestor := ObtenerGestor()
			tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
			items := crearListaModos(m.tiendaParaDev, tieneServidor)
			titulo := Icons.Server + " " + m.tiendaParaDev.Nombre
			if tieneServidor {
				titulo = Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " (servidor activo)"
			}
			m.lista = crearLista(items, titulo, m.ancho, m.alto)
			m.mensaje = ""

			return m, tea.EnableMouseCellMotion

		case "ctrl+s":

			if err := ObtenerGestor().DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = IconError(err.Error())
			} else {
				m.mensaje = Icons.Stop + " Servidor detenido"
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
						m.mensaje = IconWarning("Error enviando input")
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
		{titulo: Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
		{titulo: Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
		{titulo: Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
		{titulo: Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
	}

	if tieneServidor {
		opciones = append([]itemMenu{
			{titulo: Icons.Stop + " Detener", desc: "Parar servidor", atajo: "s"},
		}, opciones...)
	}

	return opciones
}

func (m Model) updatePopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	gestor := ObtenerGestor()
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)
	opciones := crearOpcionesPopup(tieneServidor)

	ejecutarOpcion := func(indice int) (tea.Model, tea.Cmd) {
		if indice < 0 || indice >= len(opciones) {
			return m, nil
		}

		opcion := opciones[indice]
		m.vista = VistaLogs

		switch {
		case strings.Contains(opcion.titulo, "Detener"):
			if err := gestor.DetenerServidor(m.tiendaParaDev.Nombre); err != nil {
				m.mensaje = IconError(err.Error())
			} else {
				m.mensaje = Icons.Stop + " Servidor detenido"
				m.vista = VistaMenu
				m.recrearMenuPrincipal()
			}
			return m, nil

		case strings.Contains(opcion.titulo, "Pull"):
			return m, ejecutarThemePull(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Push"):
			return m, ejecutarThemePush(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Editor"):
			return m, ejecutarAbrirEditor(m.tiendaParaDev)

		case strings.Contains(opcion.titulo, "Terminal"):
			return m, ejecutarAbrirTerminal(m.tiendaParaDev)
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
