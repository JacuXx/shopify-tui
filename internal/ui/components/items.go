package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

const shopifyGreen = "#96BF48"

// ItemMenu es un ítem genérico del menú con título, descripción y atajo.
type ItemMenu struct {
	Titulo string
	Desc   string
	Atajo  string
}

func (i ItemMenu) Title() string       { return i.Titulo }
func (i ItemMenu) Description() string { return i.Desc }
func (i ItemMenu) FilterValue() string { return i.Titulo }

// ItemTienda adapta domain.Tienda para la lista de bubbletea.
// Lleva una referencia al Manager para mostrar estado del servidor en tiempo real.
type ItemTienda struct {
	Tienda domain.Tienda
	Indice int
	Mgr    server.Manager
}

func (i ItemTienda) Title() string {
	if i.Mgr != nil && i.Mgr.TieneServidorActivo(i.Tienda.Nombre) {
		return icons.Icons.ServerOn + " " + i.Tienda.Nombre
	}
	return i.Tienda.Nombre
}

func (i ItemTienda) Description() string {
	metodo := icons.Icons.Download + " pull"
	if i.Tienda.Metodo == domain.MetodoGitClone {
		metodo = icons.Icons.Git + " git"
	}
	if i.Mgr != nil && i.Mgr.TieneServidorActivo(i.Tienda.Nombre) {
		for _, s := range i.Mgr.ObtenerServidoresActivos() {
			if s.Tienda.Nombre == i.Tienda.Nombre {
				return i.Tienda.URL + " → " + s.URL
			}
		}
	}
	return i.Tienda.URL + " [" + metodo + "]"
}

func (i ItemTienda) FilterValue() string { return i.Tienda.Nombre }

func CrearLista(items []list.Item, titulo string, ancho, alto int) list.Model {
	alturaBanner := BannerHeight(ancho)
	const margenInferior = 5
	alturaItems := alto - alturaBanner - margenInferior
	if alturaItems < 5 {
		alturaItems = 5
	}

	anchoLista := 60
	if ancho > 0 && ancho-4 > anchoLista {
		anchoLista = ancho - 4
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color(shopifyGreen)).
		BorderLeftForeground(lipgloss.Color(shopifyGreen))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color(shopifyGreen)).
		BorderLeftForeground(lipgloss.Color(shopifyGreen))

	lista := list.New(items, delegate, anchoLista, alturaItems)
	lista.Title = titulo
	lista.Styles.Title = lista.Styles.Title.
		Background(lipgloss.Color(shopifyGreen)).
		Foreground(lipgloss.Color("#FFFFFF"))
	lista.SetShowStatusBar(true)
	lista.SetFilteringEnabled(true)
	lista.SetShowPagination(true)
	return lista
}

// CrearListaTiendas convierte un slice de tiendas en items de lista.
func CrearListaTiendas(tiendas []domain.Tienda, mgr server.Manager) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = ItemTienda{Tienda: t, Indice: i + 1, Mgr: mgr}
	}
	return items
}

// CrearMenuPrincipal retorna los ítems del menú principal.
func CrearMenuPrincipal() []list.Item {
	return []list.Item{
		ItemMenu{Titulo: icons.Icons.Login + " Iniciar sesión", Desc: "Autenticarte en Shopify", Atajo: "a"},
		ItemMenu{Titulo: icons.Icons.Add + " Agregar tienda", Desc: "Registrar tienda y descargar tema", Atajo: "t"},
		ItemMenu{Titulo: icons.Icons.Server + " Desarrollo local", Desc: "Iniciar servidor", Atajo: "d"},
		ItemMenu{Titulo: icons.Icons.Logs + " Servidores activos", Desc: "Ver y administrar procesos", Atajo: "v"},
	}
}

// CrearListaMetodos retorna los métodos de descarga disponibles.
func CrearListaMetodos() []list.Item {
	return []list.Item{
		ItemMenu{Titulo: icons.Icons.Download + " Shopify Pull", Desc: "Desde Shopify directo", Atajo: "s"},
		ItemMenu{Titulo: icons.Icons.Git + " Git Clone", Desc: "Desde repositorio Git", Atajo: "g"},
	}
}

// CrearListaModos retorna las opciones de modo para una tienda dada.
func CrearListaModos(tienda domain.Tienda, tieneServidor bool) []list.Item {
	opcionesComunes := []list.Item{
		ItemMenu{Titulo: icons.Icons.Download + " Pull", Desc: "Bajar cambios del tema", Atajo: "p"},
		ItemMenu{Titulo: icons.Icons.Upload + " Push", Desc: "Subir cambios al tema", Atajo: "u"},
		ItemMenu{Titulo: icons.Icons.Editor + " Editor", Desc: "Abrir en VS Code", Atajo: "e"},
		ItemMenu{Titulo: icons.Icons.Terminal + " Terminal", Desc: "Abrir terminal aquí", Atajo: "t"},
	}

	if tieneServidor {
		items := []list.Item{
			ItemMenu{Titulo: icons.Icons.Logs + " Ver logs", Desc: "Logs en tiempo real", Atajo: "l"},
			ItemMenu{Titulo: icons.Icons.Stop + " Detener", Desc: "Parar servidor", Atajo: "s"},
		}
		return append(items, opcionesComunes...)
	}

	items := []list.Item{
		ItemMenu{Titulo: icons.Icons.Rocket + " Iniciar", Desc: "Ejecutar theme dev", Atajo: "i"},
	}
	return append(items, opcionesComunes...)
}
