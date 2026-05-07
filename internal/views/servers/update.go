package servers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

// Update maneja los eventos de la vista de servidores activos.
func Update(s State, msg tea.Msg, serverMgr server.Manager) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			servidores := serverMgr.ObtenerServidoresActivos()
			if len(servidores) == 0 {
				return s, nil
			}

			indice := s.Lista.Index()
			if indice < 0 || indice >= len(servidores) {
				indice = 0
			}

			servidor := servidores[indice]
			if err := serverMgr.DetenerServidor(servidor.Tienda.Nombre); err != nil {
				s.Mensaje = icons.IconError(err.Error())
			} else {
				s.Mensaje = icons.IconSuccess("Servidor de '" + servidor.Tienda.Nombre + "' detenido")
			}
			return s, nil

		case "S":
			serverMgr.DetenerTodos()
			s.Mensaje = icons.IconSuccess("Todos los servidores detenidos")
			return s, nil

		case "j", "down":
			servidores := serverMgr.ObtenerServidoresActivos()
			if len(servidores) > 0 {
				s.Lista, _ = s.Lista.Update(msg)
			}
			return s, nil

		case "k", "up":
			servidores := serverMgr.ObtenerServidoresActivos()
			if len(servidores) > 0 {
				s.Lista, _ = s.Lista.Update(msg)
			}
			return s, nil
		}
	}

	return s, nil
}
