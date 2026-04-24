package presentation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	styleMessage = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	styleSubtle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

func (m Model) View() string {
	switch m.vista {
	case VistaMenu:
		return m.viewMenu()
	case VistaAgregarTienda:
		return m.viewAgregarTienda()
	case VistaSeleccionarMetodo:
		return m.viewSeleccionarMetodo()
	case VistaInputGit:
		return m.viewInputGit()
	case VistaSeleccionarTienda:
		return m.viewSeleccionarTienda()
	case VistaServidores:
		return m.viewServidores()
	default:
		return m.viewMenu()
	}
}

func (m Model) viewMenu() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("🛍️ Shopify TUI\n\n"))

	sb.WriteString(m.lista.View())

	if m.mensaje != "" {
		sb.WriteString("\n\n")
		sb.WriteString(m.mensaje)
	}

	sb.WriteString("\n\n")
	sb.WriteString(styleSubtle.Render("Presiona 'q' para salir | 'ctrl+q' para forzar salida"))

	return sb.String()
}

func (m Model) viewAgregarTienda() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("➕ Agregar Tienda\n\n"))

	// Nombre
	sb.WriteString("Nombre de la tienda:\n")
	sb.WriteString(m.inputNombre.View())
	sb.WriteString("\n\n")

	// URL
	sb.WriteString("URL de la tienda:\n")
	sb.WriteString(m.inputURL.View())
	sb.WriteString("\n\n")

	if m.mensaje != "" {
		sb.WriteString(m.mensaje)
		sb.WriteString("\n\n")
	}

	sb.WriteString(styleSubtle.Render("Presiona 'Enter' para continuar | 'q' para volver"))

	return sb.String()
}

func (m Model) viewSeleccionarMetodo() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("📥 Selecciona Método de Descarga\n\n"))

	sb.WriteString(m.lista.View())

	if m.mensaje != "" {
		sb.WriteString("\n\n")
		sb.WriteString(m.mensaje)
	}

	sb.WriteString("\n\n")
	sb.WriteString(styleSubtle.Render("Presiona 'Enter' para seleccionar | 'q' para volver"))

	return sb.String()
}

func (m Model) viewInputGit() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("🌿 URL del Repositorio Git\n\n"))

	sb.WriteString("URL del repositorio:\n")
	sb.WriteString(m.inputGit.View())
	sb.WriteString("\n\n")

	if m.mensaje != "" {
		sb.WriteString(m.mensaje)
		sb.WriteString("\n\n")
	}

	sb.WriteString(styleSubtle.Render("Presiona 'Enter' para continuar | 'q' para volver"))

	return sb.String()
}

func (m Model) viewSeleccionarTienda() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("🖥️ Selecciona una Tienda\n\n"))

	if len(m.tiendas) == 0 {
		sb.WriteString(styleError.Render("No hay tiendas registradas. Agrega una primero.\n\n"))
		sb.WriteString(styleSubtle.Render("Presiona 'q' para volver"))
		return sb.String()
	}

	sb.WriteString(m.lista.View())

	if m.mensaje != "" {
		sb.WriteString("\n\n")
		sb.WriteString(m.mensaje)
	}

	sb.WriteString("\n\n")
	sb.WriteString(styleSubtle.Render("Presiona 'Enter' para seleccionar | 'q' para volver"))

	return sb.String()
}

func (m Model) viewServidores() string {
	var sb strings.Builder

	sb.WriteString(styleTitle.Render("📊 Servidores Activos\n\n"))

	activos := m.container.ServerManager.GetActiveServers()
	if len(activos) == 0 {
		sb.WriteString("No hay servidores activos.\n\n")
	} else {
		for i, srv := range activos {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, srv.Tienda.Nombre))
			sb.WriteString(fmt.Sprintf("   URL: %s\n", srv.URL))
			sb.WriteString(fmt.Sprintf("   Puerto: %d\n", srv.Puerto))
			sb.WriteString("\n")
		}
	}

	if m.mensaje != "" {
		sb.WriteString(m.mensaje)
		sb.WriteString("\n\n")
	}

	sb.WriteString(styleSubtle.Render("Presiona 'q' para volver"))

	return sb.String()
}
