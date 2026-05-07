package domain

import tea "github.com/charmbracelet/bubbletea"

// NavMsg signals a view wants to transition to a different vista.
type NavMsg struct {
	Target Vista
	Tienda *Tienda // optional context
}

func NavTo(target Vista) tea.Cmd {
	return func() tea.Msg { return NavMsg{Target: target} }
}

func NavToTienda(target Vista, tienda Tienda) tea.Cmd {
	return func() tea.Msg { return NavMsg{Target: target, Tienda: &tienda} }
}
