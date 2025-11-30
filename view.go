// view.go - Renderizado de la interfaz
// Este archivo contiene todo lo relacionado con CÓMO SE VE la aplicación
// Aquí definimos estilos (colores, bordes) y las funciones View()

package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// === ESTILOS ===
// lipgloss es como CSS pero para la terminal
// Definimos estilos reutilizables como constantes

var (
	// Estilo para títulos principales
	estiloTitulo = lipgloss.NewStyle().
			Bold(true).                            // Texto en negrita
			Foreground(lipgloss.Color("#7D56F4")). // Color morado
			MarginBottom(1)                        // Espacio abajo

	// Estilo para subtítulos
	estiloSubtitulo = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")). // Blanco
			Background(lipgloss.Color("#7D56F4")). // Fondo morado
			Padding(0, 1)                          // Padding horizontal

	// Estilo para mensajes de éxito
	estiloExito = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")). // Verde
			Bold(true)

	// Estilo para mensajes de error/advertencia
	estiloError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")). // Rojo
			Bold(true)

	// Estilo para el contenedor principal (borde redondeado)
	estiloContenedor = lipgloss.NewStyle().
				Padding(1, 2).                          // Padding interno
				Border(lipgloss.RoundedBorder()).       // Borde redondeado
				BorderForeground(lipgloss.Color("#7D56F4")) // Color del borde

	// Estilo para inputs activos (con > al inicio)
	estiloInputActivo = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	// Estilo para labels de formulario
	estiloLabel = lipgloss.NewStyle().
			Bold(true).
			MarginBottom(0)

	// Estilo para texto de ayuda (gris, abajo de la pantalla)
	estiloAyuda = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	// Estilo para texto normal
	estiloNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))
)

// View es la función principal de renderizado
// Bubbletea la llama cada vez que el estado cambia
// Debe retornar un string que representa la UI
func (m Model) View() string {
	// Dependiendo de la vista actual, mostrar diferente contenido
	switch m.vista {
	case VistaMenu:
		return m.vistaMenu()
	case VistaAgregarTienda:
		return m.vistaAgregarTienda()
	case VistaSeleccionarTienda:
		return m.vistaSeleccionarTienda()
	default:
		return m.vistaMenu()
	}
}

// vistaMenu renderiza el menú principal
func (m Model) vistaMenu() string {
	// La lista de bubbles ya tiene su propio renderizado
	s := m.lista.View()

	// Agregar mensaje si existe (éxito o error)
	if m.mensaje != "" {
		s += "\n"
		// Detectar si es éxito o error por el emoji
		if strings.HasPrefix(m.mensaje, "✅") {
			s += estiloExito.Render(m.mensaje)
		} else {
			s += estiloError.Render(m.mensaje)
		}
	}

	// Mostrar contador de tiendas guardadas
	s += "\n" + estiloAyuda.Render(
		fmt.Sprintf("📦 Tiendas guardadas: %d", len(m.tiendas)),
	)

	// Ayuda de navegación
	s += "\n" + estiloAyuda.Render("j/k: navegar • enter: seleccionar • q: salir")

	return s
}

// vistaAgregarTienda renderiza el formulario para agregar tienda
func (m Model) vistaAgregarTienda() string {
	var b strings.Builder // StringBuilder para construir el output

	// Título
	b.WriteString(estiloTitulo.Render("➕ Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	// === Campo: Nombre ===
	// Mostrar > si este campo está activo
	if m.cursorInput == 0 {
		b.WriteString(estiloInputActivo.Render("> Nombre de la tienda:"))
	} else {
		b.WriteString(estiloLabel.Render("  Nombre de la tienda:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputNombre.View()) // El input tiene su propio View()
	b.WriteString("\n\n")

	// === Campo: URL ===
	if m.cursorInput == 1 {
		b.WriteString(estiloInputActivo.Render("> URL de la tienda:"))
	} else {
		b.WriteString(estiloLabel.Render("  URL de la tienda:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputURL.View())
	b.WriteString("\n")

	// Ejemplo de URL
	b.WriteString(estiloAyuda.Render("    Ejemplo: mi-tienda.myshopify.com"))
	b.WriteString("\n\n")

	// Mostrar mensaje de error si existe
	if m.mensaje != "" {
		b.WriteString(estiloError.Render(m.mensaje))
		b.WriteString("\n\n")
	}

	// Ayuda de navegación
	b.WriteString(estiloAyuda.Render("tab: cambiar campo • enter: guardar • esc: cancelar"))

	// Envolver todo en un contenedor con borde
	return estiloContenedor.Render(b.String())
}

// vistaSeleccionarTienda renderiza la lista de tiendas para ejecutar theme dev
func (m Model) vistaSeleccionarTienda() string {
	var b strings.Builder

	// Si no hay tiendas, mostrar mensaje
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

	// Mostrar la lista de tiendas
	b.WriteString(m.lista.View())
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("enter: ejecutar theme dev • d: eliminar tienda • esc: volver"))

	return b.String()
}
