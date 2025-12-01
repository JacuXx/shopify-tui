// main.go - Punto de entrada de la aplicación
// Este es el archivo principal que inicia todo el programa

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Inicializar sistema de iconos (detecta Nerd Fonts)
	InitIcons()
	
	// Crear el programa Bubbletea
	// tea.NewProgram recibe el modelo inicial
	// Las opciones adicionales configuran comportamiento especial:
	p := tea.NewProgram(
		modeloInicial(),                 // Estado inicial de la app
		tea.WithAltScreen(),             // Usar pantalla alternativa (como vim)
		tea.WithMouseCellMotion(),       // Soporte básico de mouse
	)

	// Ejecutar el programa
	// Run() bloquea hasta que el programa termine (el usuario presione q)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error al ejecutar Shopify TUI: %v\n", err)
		os.Exit(1)
	}
}

// === NOTAS SOBRE LAS OPCIONES ===
//
// tea.WithAltScreen():
//   - Cambia a una "pantalla alternativa" de la terminal
//   - Cuando la app cierra, vuelve al estado anterior
//   - Es lo mismo que hace vim, htop, less, etc.
//
// tea.WithMouseCellMotion():
//   - Permite detectar movimiento del mouse
//   - Útil para listas scrolleables
//
// Otras opciones útiles (no las usamos ahora):
//   - tea.WithInput(r): Leer input de otro lugar
//   - tea.WithOutput(w): Escribir output a otro lugar
//   - tea.WithFilter(f): Filtrar mensajes
