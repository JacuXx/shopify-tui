package styles

import "github.com/charmbracelet/lipgloss"

var (
	Titulo = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	Exito = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	Contenedor = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))

	InputActivo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	Label = lipgloss.NewStyle().
		Bold(true).
		MarginBottom(0)

	Ayuda = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Italic(true)

	Atajo = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB800")).
		Bold(true)

	ItemNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	ItemSeleccionado = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	Desc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	Popup = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Background(lipgloss.Color("#1a1a2e"))

	PopupTitulo = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)
)
