package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

func main() {

	icons.InitIcons()

	p := tea.NewProgram(
		modeloInicial(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al ejecutar Shopify TUI: %v\n", err)
		os.Exit(1)
	}
}
