package menu

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

// State contiene el estado de la vista del menú principal.
type State struct {
	Lista list.Model
}

// NewState crea el estado inicial del menú con las dimensiones dadas.
func NewState(ancho, alto int) State {
	items := components.CrearMenuPrincipal()
	lista := components.CrearLista(items, icons.Icons.App+" Shopify TUI", ancho, alto)
	return State{Lista: lista}
}

// Rebuild recrea la lista con nuevas dimensiones.
func (s *State) Rebuild(ancho, alto int) {
	items := components.CrearMenuPrincipal()
	s.Lista = components.CrearLista(items, icons.Icons.App+" Shopify TUI", ancho, alto)
}
