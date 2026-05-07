package addstore

import (
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/JacuXx/shopify-cli/internal/domain"
)

// State contiene el estado de las vistas de agregar tienda.
type State struct {
	InputNombre      textinput.Model
	InputURL         textinput.Model
	InputGit         textinput.Model
	CursorInput      int
	TiendaTemporal   domain.Tienda
	GitURLConfirmada bool
	Mensaje          string
}

// NewState crea el estado inicial con los inputs configurados.
func NewState() State {
	inputNombre := textinput.New()
	inputNombre.Placeholder = "Mi Tienda Principal"
	inputNombre.CharLimit = 50
	inputNombre.Width = 40

	inputURL := textinput.New()
	inputURL.Placeholder = "mi-tienda"
	inputURL.CharLimit = 50
	inputURL.Width = 30

	inputGit := textinput.New()
	inputGit.Placeholder = "git@github.com:usuario/tema.git o https://..."
	inputGit.CharLimit = 200
	inputGit.Width = 50

	return State{
		InputNombre: inputNombre,
		InputURL:    inputURL,
		InputGit:    inputGit,
		CursorInput: 0,
	}
}
