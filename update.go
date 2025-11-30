// update.go - Manejo de eventos y lógica
// Este archivo contiene la función Update() que es el "cerebro" de la app
// Aquí se procesan todas las teclas y mensajes, y se actualiza el estado

package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Init se ejecuta UNA VEZ al inicio de la aplicación
// Puede retornar un comando inicial (nosotros no necesitamos ninguno)
func (m Model) Init() tea.Cmd {
	// Retornar nil significa "no hagas nada al inicio"
	return nil
}

// Update es el CORAZÓN de Bubbletea
// Se llama cada vez que algo pasa (tecla presionada, mensaje recibido, etc.)
// Recibe el mensaje y retorna: (nuevo estado, comando a ejecutar)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) es un "type switch"
	// Detecta qué tipo de mensaje llegó
	switch msg := msg.(type) {

	// === MENSAJE: Tamaño de ventana cambió ===
	case tea.WindowSizeMsg:
		m.ancho = msg.Width
		m.alto = msg.Height
		// Ajustar el tamaño de la lista
		m.lista.SetSize(msg.Width-4, msg.Height-6)
		return m, nil

	// === MENSAJE: Tecla presionada ===
	case tea.KeyMsg:
		// Manejar teclas globales (funcionan en cualquier vista)
		switch msg.String() {
		case "ctrl+c":
			// Ctrl+C siempre cierra la app
			return m, tea.Quit
		case "q":
			// q solo cierra si estamos en el menú principal
			if m.vista == VistaMenu {
				return m, tea.Quit
			}
		case "esc":
			// Esc vuelve al menú desde cualquier vista
			if m.vista != VistaMenu {
				m.vista = VistaMenu
				m.mensaje = "" // Limpiar mensajes
				m.recrearMenuPrincipal()
				return m, nil
			}
		}

	// === MENSAJE: Comando terminó exitosamente ===
	case comandoTerminadoMsg:
		m.mensaje = msg.resultado
		return m, nil

	// === MENSAJE: Error al ejecutar comando ===
	case errorMsg:
		m.mensaje = "❌ Error: " + msg.err.Error()
		return m, nil
	}

	// Si no manejamos el mensaje arriba, delegamos a la vista actual
	switch m.vista {
	case VistaMenu:
		return m.updateMenu(msg)
	case VistaAgregarTienda:
		return m.updateAgregarTienda(msg)
	case VistaSeleccionarTienda:
		return m.updateSeleccionarTienda(msg)
	}

	return m, nil
}

// updateMenu maneja eventos específicos del menú principal
func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// El usuario presionó Enter - ver qué opción seleccionó
			item, ok := m.lista.SelectedItem().(itemMenu)
			if !ok {
				return m, nil
			}

			// Ejecutar la acción correspondiente según el título
			switch item.titulo {
			case "🔐 Iniciar sesión en Shopify":
				// Ejecutar shopify auth login
				return m, ejecutarShopifyLogin()

			case "➕ Agregar tienda":
				// Cambiar a la vista del formulario
				m.vista = VistaAgregarTienda
				m.inputNombre.SetValue("")      // Limpiar inputs
				m.inputURL.SetValue("")
				m.inputNombre.Focus()           // Activar el primer input
				m.cursorInput = 0
				m.mensaje = ""
				return m, nil

			case "🚀 Ejecutar theme dev":
				// Verificar que haya tiendas
				if len(m.tiendas) == 0 {
					m.mensaje = "⚠️ No hay tiendas guardadas. Agrega una primero."
					return m, nil
				}
				// Cambiar a la vista de selección de tienda
				m.vista = VistaSeleccionarTienda
				// Crear lista con las tiendas
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

	// Actualizar el componente lista (maneja j/k, flechas, etc.)
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}

// updateAgregarTienda maneja eventos en el formulario de agregar tienda
func (m Model) updateAgregarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			// Cambiar al siguiente input
			if m.cursorInput == 0 {
				// Estamos en nombre, ir a URL
				m.cursorInput = 1
				m.inputNombre.Blur() // Desactivar nombre
				m.inputURL.Focus()   // Activar URL
			} else {
				// Estamos en URL, volver a nombre
				m.cursorInput = 0
				m.inputURL.Blur()
				m.inputNombre.Focus()
			}
			return m, nil

		case "shift+tab", "up":
			// Cambiar al input anterior
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
			// Intentar guardar la tienda
			nombre := m.inputNombre.Value()
			url := m.inputURL.Value()

			// Validar que ambos campos estén llenos
			if nombre == "" || url == "" {
				m.mensaje = "⚠️ Por favor completa ambos campos"
				return m, nil
			}

			// Crear y agregar la nueva tienda
			nuevaTienda := Tienda{
				Nombre: nombre,
				URL:    url,
			}
			m.tiendas = append(m.tiendas, nuevaTienda)

			// Guardar en el archivo JSON
			if err := guardarTiendas(m.tiendas); err != nil {
				m.mensaje = "❌ Error al guardar: " + err.Error()
				return m, nil
			}

			// Éxito - volver al menú
			m.mensaje = "✅ Tienda '" + nombre + "' guardada correctamente"
			m.vista = VistaMenu
			m.recrearMenuPrincipal()
			return m, nil
		}
	}

	// Actualizar el input activo (para que escriba los caracteres)
	var cmd tea.Cmd
	if m.cursorInput == 0 {
		m.inputNombre, cmd = m.inputNombre.Update(msg)
	} else {
		m.inputURL, cmd = m.inputURL.Update(msg)
	}
	return m, cmd
}

// updateSeleccionarTienda maneja eventos en la lista de selección de tienda
func (m Model) updateSeleccionarTienda(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Obtener la tienda seleccionada
			item, ok := m.lista.SelectedItem().(itemTienda)
			if !ok {
				return m, nil
			}
			// Ejecutar theme dev para esta tienda
			return m, ejecutarThemeDev(item.tienda.URL)

		case "d":
			// Eliminar la tienda seleccionada
			indice := m.lista.Index()
			if indice >= 0 && indice < len(m.tiendas) {
				nombreEliminada := m.tiendas[indice].Nombre
				m.tiendas = eliminarTienda(m.tiendas, indice)
				
				// Guardar cambios
				if err := guardarTiendas(m.tiendas); err != nil {
					m.mensaje = "❌ Error al eliminar: " + err.Error()
				} else {
					m.mensaje = "🗑️ Tienda '" + nombreEliminada + "' eliminada"
				}

				// Si ya no hay tiendas, volver al menú
				if len(m.tiendas) == 0 {
					m.vista = VistaMenu
					m.recrearMenuPrincipal()
					return m, nil
				}

				// Actualizar la lista
				items := crearListaTiendas(m.tiendas)
				m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
				m.lista.Title = "🚀 Selecciona una tienda"
				m.lista.SetShowStatusBar(false)
				m.lista.SetFilteringEnabled(false)
			}
			return m, nil
		}
	}

	// Actualizar el componente lista
	var cmd tea.Cmd
	m.lista, cmd = m.lista.Update(msg)
	return m, cmd
}
