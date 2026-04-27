package logs

import (
	"fmt"
	"strings"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
)

// View renderiza la vista de logs del servidor.
func View(s State, servidor *server.ServidorActivo, tiendaParaDev domain.Tienda, alto, ancho int, mensaje string) string {
	var b strings.Builder

	sepAncho := 60
	if ancho > 0 && ancho-2 < sepAncho {
		sepAncho = ancho - 2
	}

	if servidor != nil && servidor.Activo {
		b.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " " + tiendaParaDev.Nombre + " - Servidor activo"))
		b.WriteString("\n")
		b.WriteString(styles.Info.Render("   " + servidor.URL))
	} else {
		b.WriteString(styles.Error.Render(icons.Icons.Stop + " " + tiendaParaDev.Nombre + " - Servidor detenido"))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", sepAncho))
	b.WriteString("\n\n")

	if servidor != nil {
		logLines := servidor.ObtenerLogs()

		if len(logLines) == 0 {
			b.WriteString(styles.Ayuda.Render("Esperando logs del servidor..."))
			b.WriteString("\n")
		} else {
			lineasVisibles := 15
			if alto > 0 {
				lineasVisibles = alto - 12
				if lineasVisibles < 5 {
					lineasVisibles = 5
				}
			}

			inicio := s.Scroll
			fin := inicio + lineasVisibles

			if fin > len(logLines) {
				fin = len(logLines)
			}
			if inicio > len(logLines) {
				inicio = len(logLines)
			}

			for i := inicio; i < fin; i++ {
				linea := logLines[i]

				if strings.Contains(linea, "error") || strings.Contains(linea, "Error") {
					b.WriteString(styles.Error.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(styles.Exito.Render(linea))
				} else {
					b.WriteString(linea)
				}
				b.WriteString("\n")
			}

			if len(logLines) > lineasVisibles {
				b.WriteString("\n")
				porcentaje := 0
				if len(logLines)-lineasVisibles > 0 {
					porcentaje = (s.Scroll * 100) / (len(logLines) - lineasVisibles)
				}
				b.WriteString(styles.Ayuda.Render(fmt.Sprintf("Líneas %d-%d de %d (%d%%)", inicio+1, fin, len(logLines), porcentaje)))
			}
		}
	} else {
		b.WriteString(styles.Error.Render("No hay servidor activo"))
	}

	b.WriteString("\n\n")
	b.WriteString(strings.Repeat("-", sepAncho))
	b.WriteString("\n")

	if mensaje != "" {
		if strings.HasPrefix(mensaje, "✅") || strings.HasPrefix(mensaje, "🛑") {
			b.WriteString(styles.Exito.Render(mensaje))
		} else {
			b.WriteString(styles.Error.Render(mensaje))
		}
		b.WriteString("\n")
	}

	if s.ModoSeleccion {
		b.WriteString(styles.Exito.Render("MODO SELECCION ACTIVO - Selecciona texto con el mouse"))
		b.WriteString("\n")
	}

	b.WriteString(styles.Info.Render(icons.Icons.Terminal + " MODO INTERACTIVO - Las teclas se envían a Shopify CLI"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("space/m: menú | j/k: scroll | v: seleccionar | Ctrl+Q: volver"))

	return b.String()
}
