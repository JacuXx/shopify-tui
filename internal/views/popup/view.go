package popup

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
)

// View renderiza el popup de acciones rápidas.
func View(s State, tiendaParaDev domain.Tienda, serverMgr server.Manager, ancho int) string {
	tieneServidor := serverMgr.TieneServidorActivo(tiendaParaDev.Nombre)
	opciones := crearOpcionesPopup(tieneServidor)

	var popupContent strings.Builder
	popupContent.WriteString(styles.PopupTitulo.Render(icons.Icons.Rocket + " Acciones"))
	popupContent.WriteString("\n\n")

	for i, op := range opciones {
		atajo := styles.Atajo.Render("[" + strings.ToUpper(op.Atajo) + "]")
		cursor := "  "
		if i == s.Index {
			cursor = "> "
			popupContent.WriteString(cursor + atajo + " " + styles.ItemSeleccionado.Render(op.Titulo) + "\n")
		} else {
			popupContent.WriteString(cursor + atajo + " " + styles.ItemNormal.Render(op.Titulo) + "\n")
		}
	}

	popupContent.WriteString("\n")
	popupContent.WriteString(styles.Ayuda.Render("j/k navegar | l/enter ejecutar | space cerrar"))

	popupRendered := styles.Popup.Render(popupContent.String())

	servidor := serverMgr.ObtenerServidor(tiendaParaDev.Nombre)
	var header strings.Builder
	if servidor != nil && servidor.Activo {
		header.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " " + tiendaParaDev.Nombre))
		header.WriteString(" - ")
		header.WriteString(styles.Info.Render(servidor.URL))
	} else {
		header.WriteString(styles.Error.Render(icons.Icons.Stop + " " + tiendaParaDev.Nombre + " - Detenido"))
	}
	sepAncho := 50
	if ancho > 0 && ancho-2 < sepAncho {
		sepAncho = ancho - 2
	}
	header.WriteString("\n")
	header.WriteString(strings.Repeat("-", sepAncho))
	header.WriteString("\n\n")

	popupAncho := lipgloss.Width(popupRendered)
	padding := ""
	if ancho > popupAncho {
		padding = strings.Repeat(" ", (ancho-popupAncho)/2)
	}

	var resultado strings.Builder
	resultado.WriteString(header.String())

	for _, linea := range strings.Split(popupRendered, "\n") {
		resultado.WriteString(padding + linea + "\n")
	}

	return resultado.String()
}
