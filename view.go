package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

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

	estiloAtajo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB800")).
			Bold(true)

	estiloItemNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	estiloItemSeleccionado = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	estiloDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

func renderMenuConAtajos(items []itemMenu, selectedIndex int, titulo string) string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render(titulo))
	b.WriteString("\n\n")

	for i, item := range items {

		atajo := estiloAtajo.Render("[" + strings.ToUpper(item.atajo) + "]")

		var itemTitulo string
		if i == selectedIndex {
			itemTitulo = estiloItemSeleccionado.Render(item.titulo)
		} else {
			itemTitulo = estiloItemNormal.Render(item.titulo)
		}

		desc := estiloDesc.Render(item.desc)

		cursor := "  "
		if i == selectedIndex {
			cursor = "> "
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, atajo, itemTitulo))
		b.WriteString(fmt.Sprintf("      %s\n", desc))
	}

	return b.String()
}

func renderListaTiendas(tiendas []Tienda, selectedIndex int, titulo string) string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render(titulo))
	b.WriteString("\n\n")

	for i, tienda := range tiendas {

		num := estiloAtajo.Render(fmt.Sprintf("[%d]", i+1))

		var nombre string
		tieneServidor := ObtenerGestor().TieneServidorActivo(tienda.Nombre)
		if tieneServidor {
			nombre = Icons.ServerOn + " " + tienda.Nombre
		} else {
			nombre = tienda.Nombre
		}

		if i == selectedIndex {
			nombre = estiloItemSeleccionado.Render(nombre)
		} else {
			nombre = estiloItemNormal.Render(nombre)
		}

		metodo := Icons.Download + " pull"
		if tienda.Metodo == MetodoGitClone {
			metodo = Icons.Git + " git"
		}
		desc := estiloDesc.Render(tienda.URL + " [" + metodo + "]")

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
	case VistaPopup:
		return m.vistaPopup()
	default:
		return m.vistaMenu()
	}
}

func (m Model) vistaMenu() string {

	items := []itemMenu{}
	for _, item := range m.lista.Items() {
		if menuItem, ok := item.(itemMenu); ok {
			items = append(items, menuItem)
		}
	}

	s := renderMenuConAtajos(items, m.lista.Index(), Icons.App+" Shopify TUI")

	if m.mensaje != "" {
		s += "\n"
		if strings.HasPrefix(m.mensaje, "✅") || strings.HasPrefix(m.mensaje, Icons.Success) {
			s += estiloExito.Render(m.mensaje)
		} else {
			s += estiloError.Render(m.mensaje)
		}
	}

	servidoresActivos := ObtenerGestor().ContarActivos()
	s += "\n" + estiloAyuda.Render(
		fmt.Sprintf("Tiendas: %d | Servidores: %d", len(m.tiendas), servidoresActivos),
	)
	s += "\n" + estiloAyuda.Render("[A/T/D/V] j/k l/enter: seleccionar | Ctrl+Q: salir")

	return s
}

func (m Model) vistaAgregarTienda() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render("➕ Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	b.WriteString(estiloInfo.Render("Paso 1 de 2: Información básica"))
	b.WriteString("\n\n")

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

func (m Model) vistaSeleccionarMetodo() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render(Icons.Download + " Método de descarga"))
	b.WriteString("\n\n")

	b.WriteString(estiloLabel.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n")
	b.WriteString(estiloLabel.Render("URL: "))
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
		b.WriteString(estiloError.Render(m.mensaje))
		b.WriteString("\n")
	}

	b.WriteString(estiloAyuda.Render("[S]hopify Pull | [G]it Clone | l/enter | q: volver"))

	return b.String()
}

func (m Model) vistaInputGit() string {
	var b strings.Builder

	b.WriteString(estiloTitulo.Render("🔗 Clonar desde Git"))
	b.WriteString("\n\n")

	b.WriteString(estiloLabel.Render("Tienda: "))
	b.WriteString(m.tiendaTemporal.Nombre)
	b.WriteString("\n\n")

	b.WriteString(estiloInputActivo.Render("> URL del repositorio:"))
	b.WriteString("\n")
	b.WriteString("  " + m.inputGit.View())
	b.WriteString("\n\n")

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

	b.WriteString(estiloAyuda.Render("enter: clonar | q: volver"))

	return estiloContenedor.Render(b.String())
}

func (m Model) vistaSeleccionarTienda() string {
	var b strings.Builder

	if len(m.tiendas) == 0 {
		b.WriteString(estiloTitulo.Render("🚀 Ejecutar Theme Dev"))
		b.WriteString("\n\n")
		b.WriteString(estiloError.Render("No hay tiendas guardadas."))
		b.WriteString("\n")
		b.WriteString(estiloAyuda.Render("Primero agrega una tienda desde el menú principal."))
		b.WriteString("\n\n")
		b.WriteString(estiloAyuda.Render("q: volver"))
		return estiloContenedor.Render(b.String())
	}

	s := renderListaTiendas(m.tiendas, m.lista.Index(), Icons.Server+" Selecciona tienda")

	idx := m.lista.Index()
	if idx >= 0 && idx < len(m.tiendas) {
		s += estiloInfo.Render("📁 " + m.tiendas[idx].Ruta)
		s += "\n"
	}

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") || strings.HasPrefix(m.mensaje, Icons.Success) {
			s += estiloExito.Render(m.mensaje)
		} else {
			s += estiloError.Render(m.mensaje)
		}
		s += "\n"
	}

	s += estiloAyuda.Render("[1-9] l/enter: iniciar servidor | d: eliminar | q: volver")

	return s
}

func (m Model) vistaSeleccionarModo() string {
	var b strings.Builder

	gestor := ObtenerGestor()
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	if tieneServidor {
		b.WriteString(estiloExito.Render("● Servidor activo"))
		b.WriteString("\n")

		for _, s := range gestor.ObtenerServidoresActivos() {
			if s.Tienda.Nombre == m.tiendaParaDev.Nombre {
				b.WriteString(estiloInfo.Render("  " + s.URL))
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

	b.WriteString(estiloInfo.Render("📁 " + m.tiendaParaDev.Ruta))
	b.WriteString("\n")

	if m.mensaje != "" {
		if strings.HasPrefix(m.mensaje, "✅") || strings.HasPrefix(m.mensaje, Icons.Success) {
			b.WriteString(estiloExito.Render(m.mensaje))
		} else {
			b.WriteString(estiloError.Render(m.mensaje))
		}
		b.WriteString("\n")
	}

	if tieneServidor {
		b.WriteString(estiloAyuda.Render("[L]ogs [S]top [P]ull p[U]sh [E]ditor [T]erminal | q: volver"))
	} else {
		b.WriteString(estiloAyuda.Render("[I]niciar [P]ull p[U]sh [E]ditor [T]erminal | q: volver"))
	}

	return b.String()
}

func (m Model) vistaLogs() string {
	var b strings.Builder

	servidor := ObtenerGestor().ObtenerServidor(m.tiendaParaDev.Nombre)

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

	if servidor != nil {
		logs := servidor.ObtenerLogs()

		if len(logs) == 0 {
			b.WriteString(estiloAyuda.Render("Esperando logs del servidor..."))
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
					b.WriteString(estiloError.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(estiloExito.Render(linea))
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

	if m.modoSeleccion {
		b.WriteString(estiloExito.Render("✓ MODO SELECCIÓN ACTIVO - Selecciona texto con el mouse"))
		b.WriteString("\n")
	}

	b.WriteString(estiloInfo.Render(Icons.Terminal + " MODO INTERACTIVO - Las teclas se envían a Shopify CLI"))
	b.WriteString("\n")
	b.WriteString(estiloAyuda.Render("space/m: menú | j/k: scroll | v: seleccionar | Ctrl+Q: volver"))

	return b.String()
}

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
		b.WriteString(estiloAyuda.Render("q: volver"))
		return estiloContenedor.Render(b.String())
	}

	for i, servidor := range servidores {

		duracion := formatearDuracion(servidor.Iniciado)

		if i == m.lista.Index() {
			b.WriteString(estiloInputActivo.Render("> "))
		} else {
			b.WriteString("  ")
		}

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

	b.WriteString(estiloAyuda.Render("s: detener | S: detener todos | q: volver"))

	return estiloContenedor.Render(b.String())
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

var estiloPopup = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#7D56F4")).
	Padding(1, 2).
	Background(lipgloss.Color("#1a1a2e"))

var estiloPopupTitulo = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4")).
	MarginBottom(1)

func (m Model) vistaPopup() string {
	gestor := ObtenerGestor()
	tieneServidor := gestor.TieneServidorActivo(m.tiendaParaDev.Nombre)

	var opciones []itemMenu
	if tieneServidor {
		opciones = []itemMenu{
			{titulo: Icons.Stop + " Detener", desc: "Parar servidor", atajo: "s"},
			{titulo: Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
			{titulo: Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
			{titulo: Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
			{titulo: Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
		}
	} else {
		opciones = []itemMenu{
			{titulo: Icons.Download + " Pull", desc: "Bajar cambios", atajo: "p"},
			{titulo: Icons.Upload + " Push", desc: "Subir cambios", atajo: "u"},
			{titulo: Icons.Editor + " Editor", desc: "Abrir VS Code", atajo: "e"},
			{titulo: Icons.Terminal + " Terminal", desc: "Abrir terminal", atajo: "t"},
		}
	}

	var popupContent strings.Builder
	popupContent.WriteString(estiloPopupTitulo.Render(Icons.Rocket + " Acciones"))
	popupContent.WriteString("\n\n")

	for i, op := range opciones {
		atajo := estiloAtajo.Render("[" + strings.ToUpper(op.atajo) + "]")
		cursor := "  "
		if i == m.popupIndex {
			cursor = "> "
			popupContent.WriteString(cursor + atajo + " " + estiloItemSeleccionado.Render(op.titulo) + "\n")
		} else {
			popupContent.WriteString(cursor + atajo + " " + estiloItemNormal.Render(op.titulo) + "\n")
		}
	}

	popupContent.WriteString("\n")
	popupContent.WriteString(estiloAyuda.Render("j/k navegar | l/enter ejecutar | space cerrar"))

	popup := estiloPopup.Render(popupContent.String())

	servidor := gestor.ObtenerServidor(m.tiendaParaDev.Nombre)
	var header strings.Builder
	if servidor != nil && servidor.Activo {
		header.WriteString(estiloExito.Render(Icons.ServerOn + " " + m.tiendaParaDev.Nombre))
		header.WriteString(" - ")
		header.WriteString(estiloInfo.Render(servidor.URL))
	} else {
		header.WriteString(estiloError.Render(Icons.Stop + " " + m.tiendaParaDev.Nombre + " - Detenido"))
	}
	header.WriteString("\n")
	header.WriteString(strings.Repeat("─", 50))
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

	servidor := ObtenerGestor().ObtenerServidor(m.tiendaParaDev.Nombre)

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

	if servidor != nil {
		logs := servidor.ObtenerLogs()

		if len(logs) == 0 {
			b.WriteString(estiloAyuda.Render("Esperando logs del servidor..."))
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
					b.WriteString(estiloError.Render(linea))
				} else if strings.Contains(linea, "http://") || strings.Contains(linea, "https://") {
					b.WriteString(estiloExito.Render(linea))
				} else {
					b.WriteString(linea)
				}
				b.WriteString("\n")
			}
		}
	} else {
		b.WriteString(estiloError.Render("No hay servidor activo"))
	}

	return b.String()
}
