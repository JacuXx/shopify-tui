package selectstore

import (
	"fmt"
	"strings"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
)

// ViewSeleccionarTienda renderiza la vista de selección de tienda.
func ViewSeleccionarTienda(s State, ancho int) string {
	contenedor := styles.Contenedor
	if ancho > 0 {
		contenedor = contenedor.MaxWidth(ancho)
	}
	return contenedor.Render(s.Lista.View())
}

// ViewSeleccionarModo renderiza la vista de selección de modo de desarrollo.
func ViewSeleccionarModo(s State, tiendaParaDev domain.Tienda, serverMgr server.Manager) string {
	var b strings.Builder

	tieneServidor := serverMgr.TieneServidorActivo(tiendaParaDev.Nombre)

	if tieneServidor {
		b.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " Servidor activo"))
		b.WriteString("\n")

		for _, srv := range serverMgr.ObtenerServidoresActivos() {
			if srv.Tienda.Nombre == tiendaParaDev.Nombre {
				b.WriteString(styles.Info.Render("  " + srv.URL))
				b.WriteString("\n")
				break
			}
		}
		b.WriteString("\n")
	}

	items := []components.ItemMenu{}
	for _, item := range s.Lista.Items() {
		if menuItem, ok := item.(components.ItemMenu); ok {
			items = append(items, menuItem)
		}
	}

	titulo := tiendaParaDev.Nombre
	b.WriteString(renderModoConAtajos(items, s.Lista.Index(), titulo))

	b.WriteString(styles.Info.Render(icons.Icons.Folder + " " + tiendaParaDev.Ruta))
	b.WriteString("\n")

	if s.Mensaje != "" {
		if strings.HasPrefix(s.Mensaje, icons.Icons.Success) {
			b.WriteString(styles.Exito.Render(s.Mensaje))
		} else {
			b.WriteString(styles.Error.Render(s.Mensaje))
		}
		b.WriteString("\n")
	}

	if tieneServidor {
		b.WriteString(styles.Ayuda.Render("[L]ogs [S]top [P]ull p[U]sh [E]ditor [T]erminal | q: volver"))
	} else {
		b.WriteString(styles.Ayuda.Render("[I]niciar [P]ull p[U]sh [E]ditor [T]erminal | q: volver"))
	}

	return b.String()
}

func renderModoConAtajos(items []components.ItemMenu, selectedIndex int, titulo string) string {
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
