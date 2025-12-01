// view.go - Renderizado de la interfaz
// Este archivo contiene todo lo relacionado con CÓMO SE VE la aplicación

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// === ESTILOS ===
var (
	estiloTitulo = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	estiloExito = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	estiloError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	estiloContenedor = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4"))

	estiloInputActivo = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	estiloLabel = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(0)

	estiloAyuda = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	estiloInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Italic(true)
)

// View es la función principal de renderizado
func (m Model) View() string {
	switch m.vista {
	case VistaMenu:
		return m.vistaMenu()
	case VistaAgregarTienda:
		return m.vistaAgregarTienda()
	case VistaSeleccionarMetodo:
		return m.vistaSeleccionarMetodo()
	case VistaInputGit:
		return m.vistaInputGit()
	case VistaSeleccionarTienda:
		return m.vistaSeleccionarTienda()
	case VistaSeleccionarModo:
		return m.vistaSeleccionarModo()
	case VistaLogs:
		return m.vistaLogs()
	case VistaServidores:
		return m.vistaServidores()
	default:
		return m.vistaMenu()
	}
}

// vistaMenu renderiza el menú principal
func (m Model) vistaMenu() string {
	s := m.lista.View()

	if m.mensaje != "" {
		s += "\n"
		if strings.HasPrefix(m.mensaje, "✅") {
			s += estiloExito.Render(m.mensaje)
		} else {
			s += estiloError.Render(m.mensaje)
		}
	}

	// Mostrar info de tiendas y servidores
	servidoresActivos := ObtenerGestor().ContarActivos()
	s += "\n" + estiloAyuda.Render(
		fmt.Sprintf("📦 Tiendas: %d | 🟢 Servidores activos: %d", len(m.tiendas), servidoresActivos),
	)
	s += "\n" + estiloAyuda.Render("j/k: navegar • enter: seleccionar • q: salir")

	return s
}

// vistaAgregarTienda renderiza el formulario para nombre y URL
func (m Model) vistaAgregarTienda() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render("➕ Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	// Paso actual
	b.WriteString(estiloInfo.Render("Paso 1 de 2: Información básica"))
	b.WriteString("\n\n")

	// Campo: Nombre
	if m.cursorInput == 0 {
		b.WriteString(estiloInputActivo.Render("> Nombre de la tienda:"))
	} else {
		b.WriteString(estiloLabel.Render("  Nombre de la tienda:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputNombre.View())
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("    Un nombre para identificar la tienda"))
	b.WriteString("\n\n")

	// Campo: URL
	if m.cursorInput == 1 {
		b.WriteString(estiloInputActivo.Render("> URL de Shopify:"))
	} else {
		b.WriteString(estiloLabel.Render("  URL de Shopify:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputURL.View())
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("    Ejemplo: mi-tienda.myshopify.com"))
	b.WriteString("\n\n")

	if m.mensaje != "" {
		b.WriteString(estiloError.Render(m.mensaje))
		b.WriteString("\n\n")
	}

	b.WriteString(estiloAyuda.Render("tab: cambiar campo • enter: continuar • esc: cancelar"))

	return estiloContenedor.Render(b.String())
}

// vistaSeleccionarMetodo renderiza la selección de método de descarga
func (m Model) vistaSeleccionarMetodo() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render("➕ Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	// Paso actual
	b.WriteString(estiloInfo.Render("Paso 2 de 2: ¿Cómo obtener los archivos del tema?"))
	b.WriteString("\n\n")

	// Mostrar info de la tienda
	b.WriteString(estiloLabel.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n")
	b.WriteString(estiloLabel.Render("URL: "))
	b.WriteString(m.tiendaTemporal.URL)
	b.WriteString("\n\n")

	// Lista de métodos
	b.WriteString(m.lista.View())
	b.WriteString("\n")

	if m.mensaje != "" {
		b.WriteString(estiloError.Render(m.mensaje))
		b.WriteString("\n")
	}

	b.WriteString(estiloAyuda.Render("enter: seleccionar • esc: volver"))

	return b.String()
}

// vistaInputGit renderiza el input para la URL del repositorio git
func (m Model) vistaInputGit() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render("🔗 Clonar desde Git"))
	b.WriteString("\n\n")

	// Mostrar info de la tienda
	b.WriteString(estiloLabel.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n\n")

	// Input para URL de Git
	b.WriteString(estiloInputActivo.Render("> URL del repositorio:"))
	b.WriteString("\n")
	b.WriteString("  " + m.inputGit.View())
	b.WriteString("\n\n")

	// Ejemplos
	b.WriteString(estiloAyuda.Render("Ejemplos:"))
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("  SSH:   git@github.com:usuario/tema.git"))
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("  HTTPS: https://github.com/usuario/tema.git"))
	b.WriteString("\n\n")

	if m.mensaje != "" {
		b.WriteString(estiloError.Render(m.mensaje))
		b.WriteString("\n\n")
	}

	b.WriteString(estiloAyuda.Render("enter: clonar • esc: volver"))

	return estiloContenedor.Render(b.String())
}

// vistaSeleccionarTienda renderiza la lista de tiendas para theme dev
func (m Model) vistaSeleccionarTienda() string {
	var b strings.Builder

	if len(m.tiendas) == 0 {
		b.WriteString(estiloTitulo.Render("🚀 Ejecutar Theme Dev"))
		b.WriteString("\n\n")
		b.WriteString(estiloError.Render("No hay tiendas guardadas."))
		b.WriteString("\n")
		b.WriteString(estiloAyuda.Render("Primero agrega una tienda desde el menú principal."))
		b.WriteString("\n\n")
		b.WriteString(estiloAyuda.Render("esc: volver al menú"))
		return estiloContenedor.Render(b.String())
	}

	b.WriteString(m.lista.View())
	b.WriteString("\n")

	// Mostrar la ruta de la tienda seleccionada
	if item, ok := m.lista.SelectedItem().(itemTienda); ok {
		b.WriteString(estiloInfo.Render("📁 " + item.tienda.Ruta))
		b.WriteString("\n")
	}

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") {
			b.WriteString(estiloExito.Render(m.mensaje))
		} else {
			b.WriteString(estiloError.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	b.WriteString(estiloAyuda.Render("enter: seleccionar • d: eliminar • esc: volver"))

	return b.String()
}

// vistaSeleccionarModo renderiza el menú para elegir modo de servidor
func (m Model) vistaSeleccionarModo() string {
	var b strings.Builder

	// Mostrar info de la tienda
	gestor := ObtenerGestor()
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	if tieneServidor {
		b.WriteString(estiloExito.Render("🟢 Servidor activo"))
		b.WriteString("\n")
		// Mostrar URL del servidor
		for _, s := range gestor.ObtenerServidoresActivos() {
			if s.Tienda.Nombre == m.tiendaParaDev.Nombre {
				b.WriteString(estiloInfo.Render("   " + s.URL))
				b.WriteString("\n")
				break
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(m.lista.View())
	b.WriteString("\n")

	// Mostrar ruta
	b.WriteString(estiloInfo.Render("📁 " + m.tiendaParaDev.Ruta))
	b.WriteString("\n")

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") {
			b.WriteString(estiloExito.Render(m.mensaje))
		} else {
			b.WriteString(estiloError.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	b.WriteString(estiloAyuda.Render("enter: seleccionar • esc: volver"))

	return b.String()
}

// vistaLogs renderiza los logs del servidor en tiempo real
func (m Model) vistaLogs() string {
	var b strings.Builder

	servidor := ObtenerGestor().ObtenerServidor(m.tiendaParaDev.Nombre)

	// Header
	if servidor != nil && servidor.Activo {
		b.WriteString(estiloExito.Render(Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " - Servidor activo"))
		b.WriteString("\n")
		b.WriteString(estiloInfo.Render("   " + servidor.URL))
	} else {
		b.WriteString(estiloError.Render(Icons.Stop + " " + m.tiendaParaDev.Nombre + " - Servidor detenido"))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n\n")

	// Logs
	if servidor != nil {
		logs := servidor.ObtenerLogs()
		
		if len(logs) == 0 {
			b.WriteString(estiloAyuda.Render("Esperando logs del servidor..."))
			b.WriteString("\n")
		} else {
			// Determinar qué líneas mostrar (20 líneas visibles)
			lineasVisibles := 15
			if m.alto > 0 {
				lineasVisibles = m.alto - 12 // Dejar espacio para header y footer
				if lineasVisibles < 5 {
					lineasVisibles = 5
				}
			}

			inicio := m.logsScroll
			fin := inicio + lineasVisibles
			
			if fin > len(logs) {
				fin = len(logs)
			}
			if inicio > len(logs) {
				inicio = len(logs)
			}

			// Mostrar las líneas
			for i := inicio; i < fin; i++ {
				linea := logs[i]
				// Colorear según el contenido
				if strings.Contains(linea, "error") || strings.Contains(linea, "Error") {
					b.WriteString(estiloError.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(estiloExito.Render(linea))
				} else {
					b.WriteString(linea)
				}
				b.WriteString("\n")
			}

			// Indicador de scroll
			if len(logs) > lineasVisibles {
				b.WriteString("\n")
				porcentaje := 0
				if len(logs)-lineasVisibles > 0 {
					porcentaje = (m.logsScroll * 100) / (len(logs) - lineasVisibles)
				}
				b.WriteString(estiloAyuda.Render(fmt.Sprintf("Líneas %d-%d de %d (%d%%)", inicio+1, fin, len(logs), porcentaje)))
			}
		}
	} else {
		b.WriteString(estiloError.Render("No hay servidor activo"))
	}

	b.WriteString("\n\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") || strings.HasPrefix(m.mensaje, "🛑") {
			b.WriteString(estiloExito.Render(m.mensaje))
		} else {
			b.WriteString(estiloError.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	// Indicador de modo selección
	if m.modoSeleccion {
		b.WriteString(estiloExito.Render("✓ MODO SELECCIÓN ACTIVO - Selecciona texto con el mouse"))
		b.WriteString("\n")
	}

	// Ayuda - las teclas se envían a Shopify CLI excepto las de control
	b.WriteString(estiloInfo.Render(Icons.Terminal + " MODO INTERACTIVO - Las teclas se envían a Shopify CLI"))
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("Scroll: j/k, ↑/↓, PgUp/PgDn, mouse wheel"))
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("v: modo selección • g: inicio • G: final • Ctrl+Q: volver"))

	return b.String()
}

// vistaServidores renderiza la lista de servidores activos
func (m Model) vistaServidores() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render(Icons.Logs + " Servidores Activos"))
	b.WriteString("\n\n")

	servidores := ObtenerGestor().ObtenerServidoresActivos()

	if len(servidores) == 0 {
		b.WriteString(estiloAyuda.Render("No hay servidores corriendo."))
		b.WriteString("\n")
		b.WriteString(estiloAyuda.Render("Inicia uno desde '" + Icons.Rocket + " Iniciar servidor'"))
		b.WriteString("\n\n")
		b.WriteString(estiloAyuda.Render("esc: volver al menú"))
		return estiloContenedor.Render(b.String())
	}

	// Mostrar cada servidor
	for i, servidor := range servidores {
		// Calcular tiempo activo
		duracion := formatearDuracion(servidor.Iniciado)

		// Indicador de selección
		if i == m.lista.Index() {
			b.WriteString(estiloInputActivo.Render("> "))
		} else {
			b.WriteString("  ")
		}

		// Info del servidor
		b.WriteString(estiloLabel.Render(servidor.Tienda.Nombre))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("    🌐 %s\n", servidor.URL))
		b.WriteString(fmt.Sprintf("    📍 Puerto: %d | ⏱️ Activo: %s\n", servidor.Puerto, duracion))
		b.WriteString("\n")
	}

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") {
			b.WriteString(estiloExito.Render(m.mensaje))
		} else {
			b.WriteString(estiloError.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	b.WriteString(estiloAyuda.Render("s: detener servidor • S: detener todos • esc: volver"))

	return estiloContenedor.Render(b.String())
}

// formatearDuracion convierte la duración en un formato legible
func formatearDuracion(inicio time.Time) string {
	duracion := time.Since(inicio)

	if duracion < time.Minute {
		return fmt.Sprintf("%ds", int(duracion.Seconds()))
	} else if duracion < time.Hour {
		return fmt.Sprintf("%dm %ds", int(duracion.Minutes()), int(duracion.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(duracion.Hours()), int(duracion.Minutes())%60)
}
