package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/injection"
	"github.com/JacuXx/shopify-cli/presentation"
)

func main() {
	// Inicializar contenedor de dependencias
	container, err := injection.NewContainer()
	if err != nil {
		fmt.Printf("Error al inicializar contenedor: %v\n", err)
		os.Exit(1)
	}

	// Crear modelo de presentación con las dependencias inyectadas
	model, err := presentation.NewModel(container)
	if err != nil {
		fmt.Printf("Error al crear modelo: %v\n", err)
		os.Exit(1)
	}

	// Iniciar programa Bubbletea
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al ejecutar Shopify TUI: %v\n", err)
		os.Exit(1)
	}
}
