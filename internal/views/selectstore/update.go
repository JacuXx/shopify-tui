package selectstore

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/commands"
	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/store"
	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

// UpdateSeleccionarTienda maneja los eventos de la vista de selección de tienda.
// Retorna el estado actualizado, un comando, y el slice de tiendas (puede mutar por borrado).
func UpdateSeleccionarTienda(s State, msg tea.Msg, tiendas []domain.Tienda, serverMgr server.Manager, storeRepo store.Repository, ancho, alto int) (State, tea.Cmd, []domain.Tienda) {
	if s.Lista.FilterState() == list.Filtering {
		var cmd tea.Cmd
		s.Lista, cmd = s.Lista.Update(msg)
		return s, cmd, tiendas
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		var indiceSeleccionado int = -1

		if len(key) == 1 && key >= "1" && key <= "9" {
			indiceSeleccionado = int(key[0] - '1')
			if indiceSeleccionado >= len(s.Lista.VisibleItems()) {
				indiceSeleccionado = -1
			} else {
				s.Lista.Select(indiceSeleccionado)
				indiceSeleccionado = s.Lista.Index()
			}
		} else if key == "enter" || key == "l" {
			indiceSeleccionado = s.Lista.Index()
		} else if key == "d" {
			if item, ok := s.Lista.SelectedItem().(components.ItemTienda); ok {
				tienda := item.Tienda
				indiceReal := -1
				for i, t := range tiendas {
					if t.Nombre == tienda.Nombre && t.URL == tienda.URL {
						indiceReal = i
						break
					}
				}

				if indiceReal != -1 {
					nombreEliminada := tienda.Nombre
					rutaEliminada := tienda.Ruta

					if rutaEliminada != "" {
						if err := os.RemoveAll(rutaEliminada); err != nil {
							s.Mensaje = icons.IconError("Error al eliminar archivos: " + err.Error())
							return s, nil, tiendas
						}
					}

					tiendas = storeRepo.EliminarTienda(tiendas, indiceReal)

					if err := storeRepo.GuardarTiendas(tiendas); err != nil {
						s.Mensaje = icons.IconError("Error al eliminar del registro: " + err.Error())
					} else {
						s.Mensaje = icons.Icons.Delete + " Tienda '" + nombreEliminada + "' eliminada por completo"
					}

					if len(tiendas) == 0 {
						return s, domain.NavTo(domain.VistaMenu), tiendas
					}

					items := components.CrearListaTiendas(tiendas, serverMgr)
					s.Lista.SetItems(items)
				}
			}
			return s, nil, tiendas
		}

		if indiceSeleccionado >= 0 {
			if item, ok := s.Lista.SelectedItem().(components.ItemTienda); ok {
				tienda := item.Tienda

				if !storeRepo.ExisteDirectorio(tienda.Ruta) {
					s.Mensaje = icons.IconError("El directorio no existe: " + tienda.Ruta)
					return s, nil, tiendas
				}

				tieneServidor := serverMgr.TieneServidorActivo(tienda.Nombre)

				if tieneServidor {
					s.Mensaje = ""
					return s, domain.NavToTienda(domain.VistaLogs, tienda), tiendas
				}

				servidor, err := serverMgr.IniciarServidor(tienda)
				if err != nil {
					s.Mensaje = icons.IconError(err.Error())
					return s, domain.NavToTienda(domain.VistaLogs, tienda), tiendas
				}

				s.Mensaje = icons.IconSuccess("Servidor iniciado en " + servidor.URL)
				return s, domain.NavToTienda(domain.VistaLogs, tienda), tiendas
			}
		}
	}

	var cmd tea.Cmd
	s.Lista, cmd = s.Lista.Update(msg)
	return s, cmd, tiendas
}

// UpdateSeleccionarModo maneja los eventos de la vista de selección de modo de desarrollo.
func UpdateSeleccionarModo(s State, msg tea.Msg, tiendaParaDev domain.Tienda, serverMgr server.Manager, ancho, alto int) (State, tea.Cmd) {
	tieneServidor := serverMgr.TieneServidorActivo(tiendaParaDev.Nombre)

	iniciarServidor := func() (State, tea.Cmd) {
		servidor, err := serverMgr.IniciarServidor(tiendaParaDev)
		if err != nil {
			s.Mensaje = icons.IconError(err.Error())
			return s, nil
		}
		s.Mensaje = icons.IconSuccess("Servidor iniciado en " + servidor.URL)
		return s, domain.NavTo(domain.VistaLogs)
	}

	verLogs := func() (State, tea.Cmd) {
		s.Mensaje = ""
		return s, domain.NavTo(domain.VistaLogs)
	}

	detenerServidor := func() (State, tea.Cmd) {
		if err := serverMgr.DetenerServidor(tiendaParaDev.Nombre); err != nil {
			s.Mensaje = icons.IconError(err.Error())
		} else {
			s.Mensaje = icons.Icons.Stop + " Servidor detenido"
		}
		return s, domain.NavTo(domain.VistaMenu)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "i":
			if !tieneServidor {
				return iniciarServidor()
			}
		case "l":
			if tieneServidor {
				return verLogs()
			}
		case "s":
			if tieneServidor {
				return detenerServidor()
			}
		case "p":
			return s, commands.EjecutarThemePull(tiendaParaDev)
		case "u":
			return s, commands.EjecutarThemePush(tiendaParaDev)
		case "e":
			return s, commands.EjecutarAbrirEditor(tiendaParaDev)
		case "t":
			return s, commands.EjecutarAbrirTerminal(tiendaParaDev)

		case "enter":
			item, ok := s.Lista.SelectedItem().(components.ItemMenu)
			if !ok {
				return s, nil
			}

			titulo := item.Titulo

			switch {
			case strings.Contains(titulo, "Iniciar"):
				return iniciarServidor()
			case strings.Contains(titulo, "Ver logs"):
				return verLogs()
			case strings.Contains(titulo, "Detener"):
				return detenerServidor()
			case strings.Contains(titulo, "Pull"):
				return s, commands.EjecutarThemePull(tiendaParaDev)
			case strings.Contains(titulo, "Push"):
				return s, commands.EjecutarThemePush(tiendaParaDev)
			case strings.Contains(titulo, "Editor"):
				return s, commands.EjecutarAbrirEditor(tiendaParaDev)
			case strings.Contains(titulo, "Terminal"):
				return s, commands.EjecutarAbrirTerminal(tiendaParaDev)
			}
		}
	}

	var cmd tea.Cmd
	s.Lista, cmd = s.Lista.Update(msg)
	return s, cmd
}
