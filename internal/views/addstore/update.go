package addstore

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/store"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

// UpdateAgregarTienda maneja los eventos de la vista de agregar tienda (paso 1).
func UpdateAgregarTienda(s State, msg tea.Msg) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			if s.CursorInput == 0 {
				s.CursorInput = 1
				s.InputNombre.Blur()
				s.InputURL.Focus()
			} else {
				s.CursorInput = 0
				s.InputURL.Blur()
				s.InputNombre.Focus()
			}
			return s, nil

		case "shift+tab", "up":
			if s.CursorInput == 1 {
				s.CursorInput = 0
				s.InputURL.Blur()
				s.InputNombre.Focus()
			} else {
				s.CursorInput = 1
				s.InputNombre.Blur()
				s.InputURL.Focus()
			}
			return s, nil

		case "enter":
			nombre := s.InputNombre.Value()
			url := s.InputURL.Value()

			if nombre == "" || url == "" {
				s.Mensaje = icons.IconWarning("Por favor completa ambos campos")
				return s, nil
			}

			s.TiendaTemporal = domain.Tienda{
				Nombre: nombre,
				URL:    url + ".myshopify.com",
			}
			s.Mensaje = ""
			return s, domain.NavTo(domain.VistaSeleccionarMetodo)
		}
	}

	var cmd tea.Cmd
	if s.CursorInput == 0 {
		s.InputNombre, cmd = s.InputNombre.Update(msg)
	} else {
		s.InputURL, cmd = s.InputURL.Update(msg)
	}
	return s, cmd
}

// UpdateSeleccionarMetodo maneja los eventos de la vista de selección de método (atajos de teclado).
func UpdateSeleccionarMetodo(s State, msg tea.Msg, storeRepo store.Repository) (State, tea.Cmd) {
	usarShopifyPull := func() (State, tea.Cmd) {
		directorio, err := storeRepo.CrearDirectorioTienda(s.TiendaTemporal.Nombre)
		if err != nil {
			s.Mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return s, nil
		}
		s.TiendaTemporal.Metodo = domain.MetodoShopifyPull
		s.TiendaTemporal.Ruta = directorio
		return s, commands.EjecutarDescargaConExec(s.TiendaTemporal, directorio)
	}

	usarGitClone := func() (State, tea.Cmd) {
		directorio, err := storeRepo.CrearDirectorioTienda(s.TiendaTemporal.Nombre)
		if err != nil {
			s.Mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return s, nil
		}
		s.TiendaTemporal.Metodo = domain.MetodoGitClone
		s.TiendaTemporal.Ruta = directorio
		s.InputGit.SetValue("")
		s.InputGit.Focus()
		s.Mensaje = ""
		return s, domain.NavTo(domain.VistaInputGit)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "s":
			return usarShopifyPull()
		case "g":
			return usarGitClone()
		case "enter":
			// El caller debe proporcionar el ítem seleccionado vía UpdateSeleccionarMetodoConItem.
			return s, nil
		}
	}

	return s, nil
}

// UpdateSeleccionarMetodoConItem maneja el "enter" cuando el caller tiene el ítem de lista seleccionado.
func UpdateSeleccionarMetodoConItem(s State, titulo string, storeRepo store.Repository) (State, tea.Cmd) {
	usarShopifyPull := func() (State, tea.Cmd) {
		directorio, err := storeRepo.CrearDirectorioTienda(s.TiendaTemporal.Nombre)
		if err != nil {
			s.Mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return s, nil
		}
		s.TiendaTemporal.Metodo = domain.MetodoShopifyPull
		s.TiendaTemporal.Ruta = directorio
		return s, commands.EjecutarDescargaConExec(s.TiendaTemporal, directorio)
	}

	usarGitClone := func() (State, tea.Cmd) {
		directorio, err := storeRepo.CrearDirectorioTienda(s.TiendaTemporal.Nombre)
		if err != nil {
			s.Mensaje = icons.IconError("Error al crear directorio: " + err.Error())
			return s, nil
		}
		s.TiendaTemporal.Metodo = domain.MetodoGitClone
		s.TiendaTemporal.Ruta = directorio
		s.InputGit.SetValue("")
		s.InputGit.Focus()
		s.Mensaje = ""
		return s, domain.NavTo(domain.VistaInputGit)
	}

	if strings.Contains(titulo, "Shopify Pull") {
		return usarShopifyPull()
	} else if strings.Contains(titulo, "Git Clone") {
		return usarGitClone()
	}
	return s, nil
}

// UpdateInputGit maneja los eventos de la vista de input de URL Git.
func UpdateInputGit(s State, msg tea.Msg) (State, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			gitURL := strings.TrimSpace(s.InputGit.Value())
			if gitURL == "" {
				s.Mensaje = icons.IconWarning("Por favor ingresa la URL del repositorio")
				return s, nil
			}
			if !s.GitURLConfirmada {
				// Primer Enter: mostrar confirmación
				s.GitURLConfirmada = true
				s.Mensaje = icons.IconInfo("URL: " + gitURL + " — presiona Enter para confirmar")
				return s, nil
			}
			// Segundo Enter: clonar
			s.GitURLConfirmada = false
			s.TiendaTemporal.GitURL = gitURL
			return s, commands.EjecutarDescargaConExec(s.TiendaTemporal, s.TiendaTemporal.Ruta)
		case "esc", "q":
			s.GitURLConfirmada = false
		}
	}

	// Cualquier otra tecla resetea la confirmación pendiente
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() != "enter" {
			s.GitURLConfirmada = false
		}
	}

	var cmd tea.Cmd
	s.InputGit, cmd = s.InputGit.Update(msg)
	// Limpiar newlines que puedan quedar del paste
	if val := s.InputGit.Value(); strings.ContainsRune(val, '\n') || strings.ContainsRune(val, '\r') {
		s.InputGit.SetValue(strings.TrimSpace(val))
	}
	return s, cmd
}
