package popup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

func crearOpcionesPopup(tieneServidor bool) []components.ItemMenu {
	opciones := []components.ItemMenu{
		{Titulo: icons.Icons.Download + " Pull", Desc: "Bajar cambios", Atajo: "p"},
		{Titulo: icons.Icons.Upload + " Push", Desc: "Subir cambios", Atajo: "u"},
		{Titulo: icons.Icons.Editor + " Editor", Desc: "Abrir VS Code", Atajo: "e"},
		{Titulo: icons.Icons.Terminal + " Terminal", Desc: "Abrir terminal", Atajo: "t"},
	}

	if tieneServidor {
		opciones = append([]components.ItemMenu{
			{Titulo: icons.Icons.Stop + " Detener", Desc: "Parar servidor", Atajo: "s"},
		}, opciones...)
	}

	return opciones
}

// Update maneja los eventos del popup de acciones rápidas.
func Update(s State, msg tea.Msg, tiendaParaDev domain.Tienda, serverMgr server.Manager) (State, tea.Cmd) {
	tieneServidor := serverMgr.TieneServidorActivo(tiendaParaDev.Nombre)
	opciones := crearOpcionesPopup(tieneServidor)

	ejecutarOpcion := func(indice int) (State, tea.Cmd) {
		if indice < 0 || indice >= len(opciones) {
			return s, nil
		}

		opcion := opciones[indice]

		switch {
		case strings.Contains(opcion.Titulo, "Detener"):
			if err := serverMgr.DetenerServidor(tiendaParaDev.Nombre); err != nil {
				_ = icons.IconError(err.Error())
				return s, domain.NavTo(domain.VistaLogs)
			}
			return s, domain.NavTo(domain.VistaMenu)

		case strings.Contains(opcion.Titulo, "Pull"):
			return s, commands.EjecutarThemePull(tiendaParaDev)

		case strings.Contains(opcion.Titulo, "Push"):
			return s, commands.EjecutarThemePush(tiendaParaDev)

		case strings.Contains(opcion.Titulo, "Editor"):
			return s, commands.EjecutarAbrirEditor(tiendaParaDev)

		case strings.Contains(opcion.Titulo, "Terminal"):
			return s, commands.EjecutarAbrirTerminal(tiendaParaDev)
		}

		return s, domain.NavTo(domain.VistaLogs)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "esc", "q", " ", "m", "ctrl+p":
			return s, domain.NavTo(domain.VistaLogs)

		case "j", "down":
			s.Index++
			if s.Index >= len(opciones) {
				s.Index = 0
			}
			return s, nil

		case "k", "up":
			s.Index--
			if s.Index < 0 {
				s.Index = len(opciones) - 1
			}
			return s, nil

		case "enter", "l":
			return ejecutarOpcion(s.Index)

		case "s":
			if tieneServidor {
				for i, op := range opciones {
					if strings.Contains(op.Titulo, "Detener") {
						return ejecutarOpcion(i)
					}
				}
			}
		case "p":
			for i, op := range opciones {
				if strings.Contains(op.Titulo, "Pull") {
					return ejecutarOpcion(i)
				}
			}
		case "u":
			for i, op := range opciones {
				if strings.Contains(op.Titulo, "Push") {
					return ejecutarOpcion(i)
				}
			}
		case "e":
			for i, op := range opciones {
				if strings.Contains(op.Titulo, "Editor") {
					return ejecutarOpcion(i)
				}
			}
		case "t":
			for i, op := range opciones {
				if strings.Contains(op.Titulo, "Terminal") {
					return ejecutarOpcion(i)
				}
			}
		}
	}

	return s, nil
}
