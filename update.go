// update.go - Manejo de eventos y lógica
// Este archivo contiene la función Update() que procesa todas las teclas y mensajes

package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
		// Volver al menú
		m.vista = VistaMenu
		m.recrearMenuPrincipal()
		return m, nil

	case errorMsg:
		m.mensaje = "❌ Error: " + msg.err.Error()
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

			case "🚀 Ejecutar theme dev":
				if len(m.tiendas) == 0 {
					m.mensaje = "⚠️ No hay tiendas guardadas. Agrega una primero."
					return m, nil
				}
				m.vista = VistaSeleccionarTienda
				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = "🚀 Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
				m.mensaje = ""
				return m, nil

			case "❌ Salir":
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

			// Ejecutar theme dev
			return m, ejecutarThemeDev(item.tienda.URL, item.tienda.Ruta)

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
