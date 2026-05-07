package menu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
)

// Update maneja los eventos de la vista del menú principal.
func Update(s State, msg tea.Msg, tiendas []domain.Tienda) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "a":
			return s, commands.EjecutarShopifyLogin()

		case "t":
			return s, domain.NavTo(domain.VistaAgregarTienda)

		case "d":
			if len(tiendas) == 0 {
				return s, domain.NavTo(domain.VistaMenu)
			}
			return s, domain.NavTo(domain.VistaSeleccionarTienda)

		case "v":
			return s, domain.NavTo(domain.VistaServidores)

		case "enter", "l":
			item, ok := s.Lista.SelectedItem().(components.ItemMenu)
			if !ok {
				return s, nil
			}

			titulo := item.Titulo

			switch {
			case strings.Contains(titulo, "Iniciar sesión"):
				return s, commands.EjecutarShopifyLogin()

			case strings.Contains(titulo, "Agregar tienda"):
				return s, domain.NavTo(domain.VistaAgregarTienda)

			case strings.Contains(titulo, "Desarrollo local"):
				if len(tiendas) == 0 {
					return s, domain.NavTo(domain.VistaMenu)
				}
				return s, domain.NavTo(domain.VistaSeleccionarTienda)

			case strings.Contains(titulo, "Servidores activos"):
				return s, domain.NavTo(domain.VistaServidores)
			}
		}
	}

	var cmd tea.Cmd
	s.Lista, cmd = s.Lista.Update(msg)
	return s, cmd
}
