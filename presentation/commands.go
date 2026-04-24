package presentation

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/application"
	"github.com/JacuXx/shopify-cli/infrastructure/icons"
)

// addStoreResultMsg es el mensaje cuando se completa agregar una tienda
type addStoreResultMsg struct {
	err    error
	nombre string
}

// executeAddStoreCmd ejecuta el caso de uso de agregar tienda
func (m *Model) executeAddStoreCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := application.AddStoreCommand{
			Nombre: m.inputNombre.Value(),
			URL:    m.inputURL.Value(),
			Metodo: m.metodoElegido,
			GitURL: m.inputGit.Value(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tienda, err := m.container.AddStoreUseCase.Execute(ctx, cmd)
		if err != nil {
			return addStoreResultMsg{err: err}
		}

		return addStoreResultMsg{nombre: tienda.Nombre}
	}
}

// handleAddStoreResult maneja el resultado de agregar una tienda
func (m *Model) handleAddStoreResult(msg addStoreResultMsg) {
	if msg.err != nil {
		m.mensaje = icons.IconError(msg.err.Error())
		return
	}

	// Recargar tiendas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tiendas, _ := m.container.ListStoresUseCase.Execute(ctx)
	m.tiendas = tiendas

	m.mensaje = icons.IconSuccess(fmt.Sprintf("Tienda '%s' agregada correctamente", msg.nombre))
	m.vista = VistaMenu
	m.recreateMainMenu()
}
