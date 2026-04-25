package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/server"
	"github.com/JacuXx/shopify-cli/internal/store"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/version"
)

// Aliases de dominio para compatibilidad con el resto del package main.
type Vista = domain.Vista
type MetodoDescarga = domain.MetodoDescarga
type Tienda = domain.Tienda

const (
	VistaMenu             = domain.VistaMenu
	VistaAgregarTienda    = domain.VistaAgregarTienda
	VistaSeleccionarMetodo = domain.VistaSeleccionarMetodo
	VistaInputGit         = domain.VistaInputGit
	VistaSeleccionarTienda = domain.VistaSeleccionarTienda
	VistaSeleccionarModo  = domain.VistaSeleccionarModo
	VistaLogs             = domain.VistaLogs
	VistaServidores       = domain.VistaServidores
	VistaPopup            = domain.VistaPopup

	MetodoShopifyPull = domain.MetodoShopifyPull
	MetodoGitClone    = domain.MetodoGitClone
)

// Model es el estado central de la aplicación (patrón Elm Architecture).
type Model struct {
	vista Vista

	lista         list.Model
	tiendaParaDev domain.Tienda
	inputNombre   textinput.Model
	inputURL      textinput.Model
	inputGit      textinput.Model

	tiendas []domain.Tienda

	tiendaTemporal domain.Tienda
	metodoElegido  domain.MetodoDescarga

	mensaje     string
	cursorInput int
	ancho       int
	alto        int

	logsScroll       int
	modoSeleccion    bool
	popupIndex       int
	vistaAnterior    domain.Vista
	gitURLConfirmada bool

	hayActualizacion bool
	versionNueva     string

	// Dependencias inyectadas — permiten testear con mocks.
	storeRepo store.Repository
	serverMgr server.Manager
}

type itemMenu struct {
	titulo string
	desc   string
	atajo  string
}

func (i itemMenu) Title() string       { return i.titulo }
func (i itemMenu) Description() string { return i.desc }
func (i itemMenu) FilterValue() string { return i.titulo }

// itemTienda adapta domain.Tienda para la lista de bubbletea.
// Lleva una referencia al Manager para mostrar estado del servidor en tiempo real.
type itemTienda struct {
	tienda domain.Tienda
	indice int
	mgr    server.Manager
}

func (i itemTienda) Title() string {
	if i.mgr != nil && i.mgr.TieneServidorActivo(i.tienda.Nombre) {
		return icons.Icons.ServerOn + " " + i.tienda.Nombre
	}
	return i.tienda.Nombre
}

func (i itemTienda) Description() string {
	metodo := icons.Icons.Download + " pull"
	if i.tienda.Metodo == domain.MetodoGitClone {
		metodo = icons.Icons.Git + " git"
	}
	if i.mgr != nil && i.mgr.TieneServidorActivo(i.tienda.Nombre) {
		for _, s := range i.mgr.ObtenerServidoresActivos() {
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
	inputURL.Placeholder = "mi-tienda"
	inputURL.CharLimit = 50
	inputURL.Width = 30

	inputGit := textinput.New()
	inputGit.Placeholder = "git@github.com:usuario/tema.git o https://..."
	inputGit.CharLimit = 200
	inputGit.Width = 50

	storeRepo := store.NewJSONRepository()
	serverMgr := server.NewProcessManager()

	items := crearMenuPrincipal()
	lista := crearLista(items, icons.Icons.App+" Shopify TUI", 0, 0)

	tiendas, _ := storeRepo.CargarTiendas()
	hayUpdate, versionNew := version.VerificarActualizacion()

	return Model{
		vista:            domain.VistaMenu,
		lista:            lista,
		inputNombre:      inputNombre,
		inputURL:         inputURL,
		inputGit:         inputGit,
		tiendas:          tiendas,
		cursorInput:      0,
		hayActualizacion: hayUpdate,
		versionNueva:     versionNew,
		storeRepo:        storeRepo,
		serverMgr:        serverMgr,
	}
}

func crearMenuPrincipal() []list.Item {
	return []list.Item{
		itemMenu{titulo: icons.Icons.Login + " Iniciar sesión", desc: "Autenticarte en Shopify", atajo: "a"},
		itemMenu{titulo: icons.Icons.Add + " Agregar tienda", desc: "Registrar tienda y descargar tema", atajo: "t"},
		itemMenu{titulo: icons.Icons.Server + " Desarrollo local", desc: "Iniciar servidor", atajo: "d"},
		itemMenu{titulo: icons.Icons.Logs + " Servidores activos", desc: "Ver y administrar procesos", atajo: "v"},
	}
}

func crearListaTiendas(tiendas []domain.Tienda, mgr server.Manager) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = itemTienda{tienda: t, indice: i + 1, mgr: mgr}
	}
	return items
}

func crearListaMetodos() []list.Item {
	return []list.Item{
		itemMenu{titulo: icons.Icons.Download + " Shopify Pull", desc: "Desde Shopify directo", atajo: "s"},
		itemMenu{titulo: icons.Icons.Git + " Git Clone", desc: "Desde repositorio Git", atajo: "g"},
	}
}

func crearListaModos(tienda domain.Tienda, tieneServidor bool) []list.Item {
	opcionesComunes := []list.Item{
		itemMenu{titulo: icons.Icons.Download + " Pull", desc: "Bajar cambios del tema", atajo: "p"},
		itemMenu{titulo: icons.Icons.Upload + " Push", desc: "Subir cambios al tema", atajo: "u"},
		itemMenu{titulo: icons.Icons.Editor + " Editor", desc: "Abrir en VS Code", atajo: "e"},
		itemMenu{titulo: icons.Icons.Terminal + " Terminal", desc: "Abrir terminal aquí", atajo: "t"},
	}

	if tieneServidor {
		items := []list.Item{
			itemMenu{titulo: icons.Icons.Logs + " Ver logs", desc: "Logs en tiempo real", atajo: "l"},
			itemMenu{titulo: icons.Icons.Stop + " Detener", desc: "Parar servidor", atajo: "s"},
		}
		return append(items, opcionesComunes...)
	}

	items := []list.Item{
		itemMenu{titulo: icons.Icons.Rocket + " Iniciar", desc: "Ejecutar theme dev", atajo: "i"},
	}
	return append(items, opcionesComunes...)
}

func crearLista(items []list.Item, titulo string, ancho, alto int) list.Model {

	// Si hay alto definido, restamos espacio para el banner (aprox 7-8 líneas)
	// y el margen inferior.
	alturaBanner := 8
	alturaDisponible := alto - alturaBanner - 4

	alturaItems := len(items)*2 + 6
	if alto > 0 && alturaDisponible > alturaItems {
		alturaItems = alturaDisponible
	}

	if alturaItems < 10 && alto > 15 {
		alturaItems = 10
	} else if alto <= 15 {
		alturaItems = alto - 2
	}

	anchoLista := 60
	if ancho > 0 && ancho-4 > anchoLista {
		anchoLista = ancho - 4
	}

	lista := list.New(items, list.NewDefaultDelegate(), anchoLista, alturaItems)
	lista.Title = titulo
	lista.SetShowStatusBar(true)
	lista.SetFilteringEnabled(true)
	lista.SetShowPagination(true)
	return lista
}

func (m *Model) recrearMenuPrincipal() {
	items := crearMenuPrincipal()
	m.lista = crearLista(items, icons.Icons.App+" Shopify TUI", m.ancho, m.alto)
}
