package servers

import (
	"fmt"
	"strings"
	"time"

	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
)

// View renderiza la vista de servidores activos.
func View(s State, serverMgr server.Manager, ancho int) string {
	var b strings.Builder

	contenedor := styles.Contenedor
	if ancho > 0 {
		contenedor = contenedor.MaxWidth(ancho)
	}

	b.WriteString(styles.Titulo.Render(icons.Icons.Logs + " Servidores Activos"))
	b.WriteString("\n\n")

	servidores := serverMgr.ObtenerServidoresActivos()

	if len(servidores) == 0 {
		b.WriteString(styles.Ayuda.Render("No hay servidores corriendo."))
		b.WriteString("\n")
		b.WriteString(styles.Ayuda.Render("Inicia uno desde '" + icons.Icons.Rocket + " Iniciar servidor'"))
		b.WriteString("\n\n")
		b.WriteString(styles.Ayuda.Render("q: volver"))
		return contenedor.Render(b.String())
	}

	for i, servidor := range servidores {
		duracion := formatearDuracion(servidor.Iniciado)

		if i == s.Lista.Index() {
			b.WriteString(styles.InputActivo.Render("> "))
		} else {
			b.WriteString("  ")
		}

		b.WriteString(styles.Label.Render(servidor.Tienda.Nombre))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("    %s\n", servidor.URL))
		b.WriteString(fmt.Sprintf("    Puerto: %d | Activo: %s\n", servidor.Puerto, duracion))
		b.WriteString("\n")
	}

	if s.Mensaje != "" {
		if strings.HasPrefix(s.Mensaje, icons.Icons.Success) {
			b.WriteString(styles.Exito.Render(s.Mensaje))
		} else {
			b.WriteString(styles.Error.Render(s.Mensaje))
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.Ayuda.Render("s: detener | S: detener todos | q: volver"))

	return contenedor.Render(b.String())
}

func formatearDuracion(inicio time.Time) string {
	duracion := time.Since(inicio)

	if duracion < time.Minute {
		return fmt.Sprintf("%ds", int(duracion.Seconds()))
	} else if duracion < time.Hour {
		return fmt.Sprintf("%dm %ds", int(duracion.Minutes()), int(duracion.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(duracion.Hours()), int(duracion.Minutes())%60)
}
