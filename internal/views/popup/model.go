package popup

// State contiene el estado del popup de acciones rápidas.
type State struct {
	Index int
}

// NewState crea el estado inicial del popup.
func NewState() State {
	return State{}
}
