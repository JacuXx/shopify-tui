package logs

// State contiene el estado de la vista de logs.
type State struct {
	Scroll        int
	ModoSeleccion bool
}

// NewState crea el estado inicial de logs.
func NewState() State {
	return State{}
}
