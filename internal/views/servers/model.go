package servers

import (
	"github.com/charmbracelet/bubbles/list"
)

// State contiene el estado de la vista de servidores activos.
type State struct {
	Lista   list.Model
	Mensaje string
}

// NewState crea el estado inicial vacío.
func NewState() State {
	return State{
		Lista: list.New([]list.Item{}, list.NewDefaultDelegate(), 60, 10),
	}
}
