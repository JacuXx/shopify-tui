package logs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

// TickMsg es el mensaje de tick para refrescar los logs periódicamente.
type TickMsg time.Time

// tickCmd define el tick privado de 500ms.
func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// TickCmd exporta el comando de tick para que app/update.go lo inicie al navegar a logs.
func TickCmd() tea.Cmd {
	return tickCmd()
}

// Update maneja los eventos de la vista de logs.
func Update(s State, msg tea.Msg, servidor *server.ServidorActivo, tiendaParaDev domain.Tienda, serverMgr server.Manager, alto int) (State, tea.Cmd) {
	getMaxScroll := func() int {
		if servidor == nil {
			return 0
		}
		logLines := servidor.ObtenerLogs()
		maxScroll := len(logLines) - (alto - 8)
		if maxScroll < 0 {
			maxScroll = 0
		}
		return maxScroll
	}

	switch msg := msg.(type) {
	case TickMsg:
		return s, tickCmd()

	case tea.MouseMsg:
		if s.ModoSeleccion {
			return s, nil
		}

		switch msg.Type {
		case tea.MouseWheelUp:
			s.Scroll -= 3
			if s.Scroll < 0 {
				s.Scroll = 0
			}
			return s, nil
		case tea.MouseWheelDown:
			maxScroll := getMaxScroll()
			s.Scroll += 3
			if s.Scroll > maxScroll {
				s.Scroll = maxScroll
			}
			return s, nil
		}

	case tea.KeyMsg:
		key := msg.String()

		if s.ModoSeleccion {
			if key == "v" {
				s.ModoSeleccion = false
				return s, tea.EnableMouseCellMotion
			}
			return s, nil
		}

		switch key {
		case " ", "m", "ctrl+p":
			return s, domain.NavTo(domain.VistaPopup)

		case "v":
			s.ModoSeleccion = true
			return s, tea.DisableMouse

		case "ctrl+q":
			s.ModoSeleccion = false
			return s, domain.NavTo(domain.VistaSeleccionarModo)

		case "ctrl+s":
			if err := serverMgr.DetenerServidor(tiendaParaDev.Nombre); err != nil {
				_ = icons.IconError(err.Error())
			}
			return s, tickCmd()

		case "ctrl+g", "G":
			s.Scroll = getMaxScroll()
			return s, nil

		case "ctrl+t", "g":
			s.Scroll = 0
			return s, nil

		case "j", "down":
			maxScroll := getMaxScroll()
			s.Scroll++
			if s.Scroll > maxScroll {
				s.Scroll = maxScroll
			}
			return s, nil

		case "k", "up":
			s.Scroll--
			if s.Scroll < 0 {
				s.Scroll = 0
			}
			return s, nil

		case "pgdown", "ctrl+d":
			maxScroll := getMaxScroll()
			s.Scroll += 10
			if s.Scroll > maxScroll {
				s.Scroll = maxScroll
			}
			return s, nil

		case "pgup", "ctrl+u":
			s.Scroll -= 10
			if s.Scroll < 0 {
				s.Scroll = 0
			}
			return s, nil

		default:
			if servidor != nil && servidor.Activo {
				var input string
				switch key {
				case "enter":
					input = "\n"
				case "space":
					input = " "
				case "tab":
					input = "\t"
				case "backspace":
					input = "\b"
				default:
					if len(key) == 1 {
						input = key
					}
				}

				if input != "" {
					servidor.EnviarInput(input)
				}
			}
			return s, nil
		}
	}

	return s, nil
}
