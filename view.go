package main

import (
	"fmt"
	"strings"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
	"github.com/JacuXx/shopify-cli/internal/version"
)

// Estilos centralizados en internal/ui/styles

func renderBanner() string {
	myFigure := figure.NewFigure("SHOPIFY TUI", "slant", true)
	asciiArt := myFigure.String()

	verdeShopify := lipgloss.NewStyle().Foreground(lipgloss.Color("#95BF47")).Bold(true)

	lineas := strings.Split(asciiArt, "\n")
	var banner strings.Builder
	for _, linea := range lineas {
		if strings.TrimSpace(linea) == "" {
			continue
		}
		banner.WriteString(verdeShopify.Render(linea))
		banner.WriteString("\n")
	}
	return banner.String()
}

func renderMenuConAtajos(items []itemMenu, selectedIndex int, titulo string) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(titulo))
	b.WriteString("\n\n")

	for i, item := range items {

		atajo := styles.Atajo.Render("[" + strings.ToUpper(item.atajo) + "]")

		var itemTitulo string
		if i == selectedIndex {
			itemTitulo = styles.ItemSeleccionado.Render(item.titulo)
		} else {
			itemTitulo = styles.ItemNormal.Render(item.titulo)
		}

		desc := styles.Desc.Render(item.desc)

		cursor := "  "
		if i == selectedIndex {
			cursor = "> "
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, atajo, itemTitulo))
		b.WriteString(fmt.Sprintf("      %s\n", desc))
	}

	return b.String()
}

func renderListaTiendas(tiendas []Tienda, selectedIndex int, titulo string, mgr interface {
	TieneServidorActivo(string) bool
}) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(titulo))
	b.WriteString("\n\n")

	for i, tienda := range tiendas {

		num := styles.Atajo.Render(fmt.Sprintf("[%d]", i+1))

		var nombre string
		tieneServidor := mgr.TieneServidorActivo(tienda.Nombre)
		if tieneServidor {
			nombre = icons.Icons.ServerOn + " " + tienda.Nombre
		} else {
			nombre = tienda.Nombre
		}

		if i == selectedIndex {
			nombre = styles.ItemSeleccionado.Render(nombre)
		} else {
			nombre = styles.ItemNormal.Render(nombre)
		}

		metodo := icons.Icons.Download + " pull"
		if tienda.Metodo == MetodoGitClone {
			metodo = icons.Icons.Git + " git"
		}
		desc := styles.Desc.Render(tienda.URL + " [" + metodo + "]")

		cursor := "  "
		if i == selectedIndex {
			cursor = "> "
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, num, nombre))
		b.WriteString(fmt.Sprintf("      %s\n", desc))
	}

	return b.String()
}

func (m Model) View() string {
	banner := renderBanner()
	var contenido string

	switch m.vista {
	case VistaMenu:
		contenido = m.vistaMenu()
	case VistaAgregarTienda:
		contenido = m.vistaAgregarTienda()
	case VistaSeleccionarMetodo:
		contenido = m.vistaSeleccionarMetodo()
	case VistaInputGit:
		contenido = m.vistaInputGit()
	case VistaSeleccionarTienda:
		contenido = m.vistaSeleccionarTienda()
	case VistaSeleccionarModo:
		contenido = m.vistaSeleccionarModo()
	case VistaLogs:
		contenido = m.vistaLogs()
	case VistaServidores:
		contenido = m.vistaServidores()
	case VistaPopup:
		contenido = m.vistaPopup()
	default:
		contenido = m.vistaMenu()
	}

	if m.vista == VistaMenu {
		return contenido
	}

	return banner + "\n" + contenido
}

func (m Model) vistaMenu() string {

	items := []itemMenu{}
	for _, item := range m.lista.Items() {
		if menuItem, ok := item.(itemMenu); ok {
			items = append(items, menuItem)
		}
	}

	banner := renderBanner()
	menu := renderMenuConAtajos(items, m.lista.Index(), icons.Icons.App+" Configuración")

	s := banner + "\n" + menu

	if m.hayActualizacion {
		estiloAviso := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
		estiloVersion := lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
		estiloComando := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		s += "\n\n" + estiloAviso.Render("Nueva version disponible: ") + estiloVersion.Render(m.versionNueva) + estiloAviso.Render(" (actual: "+version.Version+")")
		s += "\n" + estiloAviso.Render("📦 Actualiza: ") + estiloComando.Render("npm update -g shopify-cli-tui")
	}

	if m.mensaje != "" {
		s += "\n"
		if strings.HasPrefix(m.mensaje, "âœ…") || strings.HasPrefix(m.mensaje, icons.Icons.Success) {
			s += styles.Exito.Render(m.mensaje)
		} else {
			s += styles.Error.Render(m.mensaje)
		}
	}

	servidoresActivos := m.serverMgr.ContarActivos()
	s += "\n" + styles.Ayuda.Render(
		fmt.Sprintf("Tiendas: %d | Servidores: %d", len(m.tiendas), servidoresActivos),
	)
	s += "\n" + styles.Ayuda.Render("[A/T/D/V] j/k l/enter: seleccionar | Ctrl+Q: salir")

	return s
}

func (m Model) vistaAgregarTienda() string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render("âž• Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	b.WriteString(styles.Info.Render("Paso 1 de 2: Información básica"))
	b.WriteString("\n\n")

	if m.cursorInput == 0 {
		b.WriteString(styles.InputActivo.Render("> Nombre de la tienda:"))
	} else {
		b.WriteString(styles.Label.Render("  Nombre de la tienda:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputNombre.View())
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("    Un nombre para identificar la tienda"))
	b.WriteString("\n\n")

	if m.cursorInput == 1 {
		b.WriteString(styles.InputActivo.Render("> URL de Shopify:"))
	} else {
		b.WriteString(styles.Label.Render("  URL de Shopify:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + m.inputURL.View() + lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Render(".myshopify.com"))
	b.WriteString("\n\n")

	if m.mensaje != "" {
		b.WriteString(styles.Error.Render(m.mensaje))
		b.WriteString("\n\n")
	}

	b.WriteString(styles.Ayuda.Render("tab: cambiar campo â€¢ enter: continuar â€¢ esc: cancelar"))

	return styles.Contenedor.Render(b.String())
}

func (m Model) vistaSeleccionarMetodo() string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(icons.Icons.Download + " Método de descarga"))
	b.WriteString("\n\n")

	b.WriteString(styles.Label.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n")
	b.WriteString(styles.Label.Render("URL: "))
	b.WriteString(m.tiendaTemporal.URL)
	b.WriteString("\n\n")

	items := []itemMenu{}
	for _, item := range m.lista.Items() {
		if menuItem, ok := item.(itemMenu); ok {
			items = append(items, menuItem)
		}
	}

	b.WriteString(renderMenuConAtajos(items, m.lista.Index(), "Elige método"))
	b.WriteString("\n")

	if m.mensaje != "" {
		b.WriteString(styles.Error.Render(m.mensaje))
		b.WriteString("\n")
	}

	b.WriteString(styles.Ayuda.Render("[S]hopify Pull | [G]it Clone | l/enter | q: volver"))

	return b.String()
}

func (m Model) vistaInputGit() string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(icons.Icons.Git + " Clonar desde Git"))
	b.WriteString("\n\n")

	b.WriteString(styles.Label.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n\n")

	b.WriteString(styles.InputActivo.Render("> URL del repositorio:"))
	b.WriteString("\n")
	b.WriteString("  " + m.inputGit.View())
	b.WriteString("\n\n")

	b.WriteString(styles.Ayuda.Render("Ejemplos:"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  SSH:    git@github.com:usuario/tema.git"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  HTTPS:  https://github.com/usuario/tema.git"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  Privado: https://TOKEN@github.com/usuario/tema.git"))
	b.WriteString("\n\n")

	if m.mensaje != "" {
		b.WriteString(styles.Error.Render(m.mensaje))
		b.WriteString("\n\n")
	}

	if m.gitURLConfirmada {
		b.WriteString(styles.Info.Render("enter: confirmar clone | cualquier tecla: cancelar"))
	} else {
		b.WriteString(styles.Ayuda.Render("enter: continuar | q: volver"))
	}

	return styles.Contenedor.Render(b.String())
}

func (m Model) vistaSeleccionarTienda() string {
	if len(m.tiendas) == 0 {
		var b strings.Builder
		b.WriteString(styles.Titulo.Render("ðŸš€ Ejecutar Theme Dev"))
		b.WriteString("\n\n")
		b.WriteString(styles.Error.Render("No hay tiendas guardadas."))
		b.WriteString("\n")
		b.WriteString(styles.Ayuda.Render("Primero agrega una tienda desde el menú principal."))
		b.WriteString("\n\n")
		b.WriteString(styles.Ayuda.Render("q: volver"))
		return styles.Contenedor.Render(b.String())
	}

	// Usamos la m.lista.View() directamente para tener buscador y paginación
	return styles.Contenedor.Render(m.lista.View())
}

func (m Model) vistaSeleccionarModo() string {
	var b strings.Builder

	gestor := m.serverMgr
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	if tieneServidor {
		b.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " Servidor activo"))
		b.WriteString("\n")

		for _, s := range gestor.ObtenerServidoresActivos() {
			if s.Tienda.Nombre == m.tiendaParaDev.Nombre {
				b.WriteString(styles.Info.Render("  " + s.URL))
				b.WriteString("\n")
				break
			}
		}
		b.WriteString("\n")
	}

	items := []itemMenu{}
	for _, item := range m.lista.Items() {
		if menuItem, ok := item.(itemMenu); ok {
			items = append(items, menuItem)
		}
	}

	titulo := m.tiendaParaDev.Nombre
	b.WriteString(renderMenuConAtajos(items, m.lista.Index(), titulo))

	b.WriteString(styles.Info.Render(icons.Icons.Folder + " " + m.tiendaParaDev.Ruta))
	b.WriteString("\n")

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "âœ…") || strings.HasPrefix(m.mensaje, icons.Icons.Success) {
			b.WriteString(styles.Exito.Render(m.mensaje))
		} else {
			b.WriteString(styles.Error.Render(m.mensaje))
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

func (m Model) vistaLogs() string {
	var b strings.Builder

	servidor := m.serverMgr.ObtenerServidor(m.tiendaParaDev.Nombre)

	if servidor != nil && servidor.Activo {
		b.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " - Servidor activo"))
		b.WriteString("\n")
		b.WriteString(styles.Info.Render("   " + servidor.URL))
	} else {
		b.WriteString(styles.Error.Render(icons.Icons.Stop + " " + m.tiendaParaDev.Nombre + " - Servidor detenido"))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 60))
	b.WriteString("\n\n")

	if servidor != nil {
		logs := servidor.ObtenerLogs()

		if len(logs) == 0 {
			b.WriteString(styles.Ayuda.Render("Esperando logs del servidor..."))
			b.WriteString("\n")
		} else {

			lineasVisibles := 15
			if m.alto > 0 {
				lineasVisibles = m.alto - 12
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

			for i := inicio; i < fin; i++ {
				linea := logs[i]

				if strings.Contains(linea, "error") || strings.Contains(linea, "Error") {
					b.WriteString(styles.Error.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(styles.Exito.Render(linea))
				} else {
					b.WriteString(linea)
				}
				b.WriteString("\n")
			}

			if len(logs) > lineasVisibles {
				b.WriteString("\n")
				porcentaje := 0
				if len(logs)-lineasVisibles > 0 {
					porcentaje = (m.logsScroll * 100) / (len(logs) - lineasVisibles)
				}
				b.WriteString(styles.Ayuda.Render(fmt.Sprintf("Líneas %d-%d de %d (%d%%)", inicio+1, fin, len(logs), porcentaje)))
			}
		}
	} else {
		b.WriteString(styles.Error.Render("No hay servidor activo"))
	}

	b.WriteString("\n\n")
	b.WriteString(strings.Repeat("-", 60))
	b.WriteString("\n")

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "âœ…") || strings.HasPrefix(m.mensaje, "ðŸ›‘") {
			b.WriteString(styles.Exito.Render(m.mensaje))
		} else {
			b.WriteString(styles.Error.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	if m.modoSeleccion {
		b.WriteString(styles.Exito.Render("MODO SELECCION ACTIVO - Selecciona texto con el mouse"))
		b.WriteString("\n")
	}

	b.WriteString(styles.Info.Render(icons.Icons.Terminal + " MODO INTERACTIVO - Las teclas se envían a Shopify CLI"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("space/m: menú | j/k: scroll | v: seleccionar | Ctrl+Q: volver"))

	return b.String()
}

func (m Model) vistaServidores() string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(icons.Icons.Logs + " Servidores Activos"))
	b.WriteString("\n\n")

	servidores := m.serverMgr.ObtenerServidoresActivos()

	if len(servidores) == 0 {
		b.WriteString(styles.Ayuda.Render("No hay servidores corriendo."))
		b.WriteString("\n")
		b.WriteString(styles.Ayuda.Render("Inicia uno desde '" + icons.Icons.Rocket + " Iniciar servidor'"))
		b.WriteString("\n\n")
		b.WriteString(styles.Ayuda.Render("q: volver"))
		return styles.Contenedor.Render(b.String())
	}

	for i, servidor := range servidores {

		duracion := formatearDuracion(servidor.Iniciado)

		if i == m.lista.Index() {
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

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "âœ…") {
			b.WriteString(styles.Exito.Render(m.mensaje))
		} else {
			b.WriteString(styles.Error.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.Ayuda.Render("s: detener | S: detener todos | q: volver"))

	return styles.Contenedor.Render(b.String())
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

func (m Model) vistaPopup() string {
	gestor := m.serverMgr
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	var opciones []itemMenu
	if tieneServidor {
		opciones = []itemMenu{
			{titulo: icons.Icons.Stop + " Detener", desc: "Parar servidor", atajo: "s"},
			{titulo: icons.Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
			{titulo: icons.Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
			{titulo: icons.Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
			{titulo: icons.Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
		}
	} else {
		opciones = []itemMenu{
			{titulo: icons.Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
			{titulo: icons.Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
			{titulo: icons.Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
			{titulo: icons.Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
		}
	}

	var popupContent strings.Builder
	popupContent.WriteString(styles.PopupTitulo.Render(icons.Icons.Rocket + " Acciones"))
	popupContent.WriteString("\n\n")

	for i, op := range opciones {
		atajo := styles.Atajo.Render("[" + strings.ToUpper(op.atajo) + "]")
		cursor := "  "
		if i == m.popupIndex {
			cursor = "> "
			popupContent.WriteString(cursor + atajo + " " + styles.ItemSeleccionado.Render(op.titulo) + "\n")
		} else {
			popupContent.WriteString(cursor + atajo + " " + styles.ItemNormal.Render(op.titulo) + "\n")
		}
	}

	popupContent.WriteString("\n")
	popupContent.WriteString(styles.Ayuda.Render("j/k navegar | l/enter ejecutar | space cerrar"))

	popup := styles.Popup.Render(popupContent.String())

	servidor := gestor.ObtenerServidor(m.tiendaParaDev.Nombre)
	var header strings.Builder
	if servidor != nil && servidor.Activo {
		header.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre))
		header.WriteString(" - ")
		header.WriteString(styles.Info.Render(servidor.URL))
	} else {
		header.WriteString(styles.Error.Render(icons.Icons.Stop + " " + m.tiendaParaDev.Nombre + " - Detenido"))
	}
	header.WriteString("\n")
	header.WriteString(strings.Repeat("-", 50))
	header.WriteString("\n\n")

	popupAncho := lipgloss.Width(popup)
	padding := ""
	if m.ancho > popupAncho {
		padding = strings.Repeat(" ", (m.ancho-popupAncho)/2)
	}

	var resultado strings.Builder
	resultado.WriteString(header.String())

	for _, linea := range strings.Split(popup, "\n") {
		resultado.WriteString(padding + linea + "\n")
	}

	return resultado.String()
}

func (m Model) vistaLogsBase() string {
	var b strings.Builder

	servidor := m.serverMgr.ObtenerServidor(m.tiendaParaDev.Nombre)

	if servidor != nil && servidor.Activo {
		b.WriteString(styles.Exito.Render(icons.Icons.ServerOn + " " + m.tiendaParaDev.Nombre + " - Servidor activo"))
		b.WriteString("\n")
		b.WriteString(styles.Info.Render("   " + servidor.URL))
	} else {
		b.WriteString(styles.Error.Render(icons.Icons.Stop + " " + m.tiendaParaDev.Nombre + " - Servidor detenido"))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("-", 60))
	b.WriteString("\n\n")

	if servidor != nil {
		logs := servidor.ObtenerLogs()

		if len(logs) == 0 {
			b.WriteString(styles.Ayuda.Render("Esperando logs del servidor..."))
			b.WriteString("\n")
		} else {
			lineasVisibles := 15
			if m.alto > 0 {
				lineasVisibles = m.alto - 12
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

			for i := inicio; i < fin; i++ {
				linea := logs[i]
				if strings.Contains(linea, "error") || strings.Contains(linea, "Error") {
					b.WriteString(styles.Error.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(styles.Exito.Render(linea))
				} else {
					b.WriteString(linea)
				}
				b.WriteString("\n")
			}
		}
	} else {
		b.WriteString(styles.Error.Render("No hay servidor activo"))
	}

	return b.String()
}
