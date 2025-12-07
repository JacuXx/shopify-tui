package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Vista int

const (
	VistaMenu Vista = iota
	VistaAgregarTienda
	VistaSeleccionarMetodo
	VistaInputGit
	VistaSeleccionarTienda
	VistaSeleccionarModo
	VistaLogs
	VistaServidores
)

type MetodoDescarga int

const (
	MetodoShopifyPull MetodoDescarga = iota
	MetodoGitClone
)

type Tienda struct {
	Nombre string         `json:"nombre"`
	URL    string         `json:"url"`
	Ruta   string         `json:"ruta"`
	Metodo MetodoDescarga `json:"metodo"`
	GitURL string         `json:"git_url,omitempty"`
}

type Model struct {
	vista Vista

	lista         list.Model
	tiendaParaDev Tienda
	inputNombre   textinput.Model
	inputURL      textinput.Model
	inputGit      textinput.Model

	tiendas []Tienda

	tiendaTemporal Tienda
	metodoElegido  MetodoDescarga

	mensaje     string
	cursorInput int
	ancho       int
	alto        int

	logsScroll    int
	modoSeleccion bool
}

type itemMenu struct {
	titulo string
	desc   string
	atajo  string
}

func (i itemMenu) Title() string       { return i.titulo }
func (i itemMenu) Description() string { return i.desc }
func (i itemMenu) FilterValue() string { return i.titulo }

type itemTienda struct {
	tienda Tienda
	indice int
}

func (i itemTienda) Title() string {

	if ObtenerGestor().TieneServidorActivo(i.tienda.Nombre) {
		return Icons.ServerOn + " " + i.tienda.Nombre
	}
	return i.tienda.Nombre
}
func (i itemTienda) Description() string {
	metodo := Icons.Download + " pull"
	if i.tienda.Metodo == MetodoGitClone {
		metodo = Icons.Git + " git"
	}

	if ObtenerGestor().TieneServidorActivo(i.tienda.Nombre) {
		servidores := ObtenerGestor().ObtenerServidoresActivos()
		for _, s := range servidores {
			if s.Tienda.Nombre == i.tienda.Nombre {
				return i.tienda.URL + " → " + s.URL
			}
		}
	}
	return i.tienda.URL + " [" + metodo + "]"
}
func (i itemTienda) FilterValue() string { return i.tienda.Nombre }

func modeloInicial() Model {

	inputNombre := textinput.New()
	inputNombre.Placeholder = "Mi Tienda Principal"
	inputNombre.CharLimit = 50
	inputNombre.Width = 40

	inputURL := textinput.New()
	inputURL.Placeholder = "mi-tienda.myshopify.com"
	inputURL.CharLimit = 100
	inputURL.Width = 40

	inputGit := textinput.New()
	inputGit.Placeholder = "git@github.com:usuario/tema.git o https://..."
	inputGit.CharLimit = 200
	inputGit.Width = 50

	items := crearMenuPrincipal()

	lista := crearLista(items, Icons.App+" Shopify TUI", 0, 0)

	tiendas, _ := cargarTiendas()

	return Model{
		vista:       VistaMenu,
		lista:       lista,
		inputNombre: inputNombre,
		inputURL:    inputURL,
		inputGit:    inputGit,
		tiendas:     tiendas,
		cursorInput: 0,
	}
}

func crearMenuPrincipal() []list.Item {
	return []list.Item{
		itemMenu{
			titulo: Icons.Login + " Iniciar sesión",
			desc:   "Autenticarte en Shopify",
			atajo:  "a",
		},
		itemMenu{
			titulo: Icons.Add + " Agregar tienda",
			desc:   "Registrar tienda y descargar tema",
			atajo:  "t",
		},
		itemMenu{
			titulo: Icons.Server + " Desarrollo local",
			desc:   "Iniciar servidor",
			atajo:  "d",
		},
		itemMenu{
			titulo: Icons.Logs + " Servidores activos",
			desc:   "Ver y administrar procesos",
			atajo:  "v",
		},
	}
}

func crearListaTiendas(tiendas []Tienda) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = itemTienda{tienda: t, indice: i + 1}
	}
	return items
}

func crearListaMetodos() []list.Item {
	return []list.Item{
		itemMenu{
			titulo: Icons.Download + " Shopify Pull",
			desc:   "Desde Shopify directo",
			atajo:  "s",
		},
		itemMenu{
			titulo: Icons.Git + " Git Clone",
			desc:   "Desde repositorio Git",
			atajo:  "g",
		},
	}
}

func crearListaModos(tienda Tienda, tieneServidor bool) []list.Item {

	opcionesComunes := []list.Item{
		itemMenu{
			titulo: Icons.Download + " Pull",
			desc:   "Bajar cambios del tema",
			atajo:  "p",
		},
		itemMenu{
			titulo: Icons.Upload + " Push",
			desc:   "Subir cambios al tema",
			atajo:  "u",
		},
		itemMenu{
			titulo: Icons.Editor + " Editor",
			desc:   "Abrir en VS Code",
			atajo:  "e",
		},
		itemMenu{
			titulo: Icons.Terminal + " Terminal",
			desc:   "Abrir terminal aquí",
			atajo:  "t",
		},
	}

	if tieneServidor {

		items := []list.Item{
			itemMenu{
				titulo: Icons.Logs + " Ver logs",
				desc:   "Logs en tiempo real",
				atajo:  "l",
			},
			itemMenu{
				titulo: Icons.Stop + " Detener",
				desc:   "Parar servidor",
				atajo:  "s",
			},
		}
		return append(items, opcionesComunes...)
	}

	items := []list.Item{
		itemMenu{
			titulo: Icons.Rocket + " Iniciar",
			desc:   "Ejecutar theme dev",
			atajo:  "i",
		},
	}
	return append(items, opcionesComunes...)
}

func crearLista(items []list.Item, titulo string, ancho, alto int) list.Model {

	alturaItems := len(items)*2 + 4
	if alto > 0 && alto-6 > alturaItems {
		alturaItems = alto - 6
	}
	if alturaItems < 10 {
		alturaItems = 10
	}

	anchoLista := 55
	if ancho > 0 && ancho-4 > anchoLista {
		anchoLista = ancho - 4
	}

	lista := list.New(items, list.NewDefaultDelegate(), anchoLista, alturaItems)
	lista.Title = titulo
	lista.SetShowStatusBar(false)
	lista.SetFilteringEnabled(false)
	lista.SetShowPagination(false)
	return lista
}

func (m *Model) recrearMenuPrincipal() {
	items := crearMenuPrincipal()
	m.lista = crearLista(items, Icons.App+" Shopify TUI", m.ancho, m.alto)
}
