package menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
	"github.com/JacuXx/shopify-cli/internal/version"
)

// View renderiza la vista del menú principal completa (incluyendo banner).
func View(s State, tiendas []domain.Tienda, serverMgr server.Manager, hayActualizacion bool, versionNueva, mensaje string) string {
	items := []components.ItemMenu{}
	for _, item := range s.Lista.Items() {
		if menuItem, ok := item.(components.ItemMenu); ok {
			items = append(items, menuItem)
		}
	}

	banner := components.RenderBanner()
	menuStr := renderMenuConAtajos(items, s.Lista.Index(), icons.Icons.App+" Configuración")

	result := banner + "\n" + menuStr

	if hayActualizacion {
		estiloAviso := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
		estiloVersion := lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
		estiloComando := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		result += "\n\n" + estiloAviso.Render("Nueva version disponible: ") + estiloVersion.Render(versionNueva) + estiloAviso.Render(" (actual: "+version.Version+")")
		result += "\n" + estiloAviso.Render("📦 Actualiza: ") + estiloComando.Render("bun add -g @jacuxx/shopify-tui --registry https://npm.jsr.io")
	}

	if mensaje != "" {
		result += "\n"
		if strings.HasPrefix(mensaje, "✅") || strings.HasPrefix(mensaje, icons.Icons.Success) {
			result += styles.Exito.Render(mensaje)
		} else {
			result += styles.Error.Render(mensaje)
		}
	}

	servidoresActivos := serverMgr.ContarActivos()
	result += "\n" + styles.Ayuda.Render(
		fmt.Sprintf("Tiendas: %d | Servidores: %d", len(tiendas), servidoresActivos),
	)
	result += "\n" + styles.Ayuda.Render("[A/T/D/V] j/k l/enter: seleccionar | Ctrl+Q: salir")

	return result
}

func renderMenuConAtajos(items []components.ItemMenu, selectedIndex int, titulo string) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(titulo))
	b.WriteString("\n\n")

	for i, item := range items {
		atajo := styles.Atajo.Render("[" + strings.ToUpper(item.Atajo) + "]")

		var itemTitulo string
		if i == selectedIndex {
			itemTitulo = styles.ItemSeleccionado.Render(item.Titulo)
		} else {
			itemTitulo = styles.ItemNormal.Render(item.Titulo)
		}

		desc := styles.Desc.Render(item.Desc)

		cursor := "  "
		if i == selectedIndex {
			cursor = "> "
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, atajo, itemTitulo))
		b.WriteString(fmt.Sprintf("      %s\n", desc))
	}

	return b.String()
}
