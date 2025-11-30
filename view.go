// view.go - Renderizado de la interfaz
// Este archivo contiene todo lo relacionado con CÓMO SE VE la aplicación

package main

import (
	"fmt"
	"strings"

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

	s += "\n" + estiloAyuda.Render(
		fmt.Sprintf("📦 Tiendas guardadas: %d", len(m.tiendas)),
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

	b.WriteString(estiloAyuda.Render("enter: ejecutar theme dev • d: eliminar • esc: volver"))

	return b.String()
}
