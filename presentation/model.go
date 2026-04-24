package presentation

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/JacuXx/shopify-cli/domain"
	"github.com/JacuXx/shopify-cli/infrastructure/icons"
	"github.com/JacuXx/shopify-cli/injection"
)

// Vista representa la vista actual de la aplicación
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
	VistaPopup
)

// Model es el modelo de la aplicación (Elm Architecture)
type Model struct {
	// Dependencias inyectadas
	container *injection.Container

	// Estado de UI
	vista Vista

	// Componentes Bubbletea
	lista         list.Model
	tiendaParaDev domain.Tienda
	inputNombre   textinput.Model
	inputURL      textinput.Model
	inputGit      textinput.Model

	// Estado de negocio
	tiendas []domain.Tienda

	// Temporal
	tiendaTemporal domain.Tienda
	metodoElegido  int // 0 = shopify-pull, 1 = git-clone

	// UI
	mensaje       string
	cursorInput   int
	ancho         int
	alto          int
	logsScroll    int
	modoSeleccion bool
	popupIndex    int
	vistaAnterior Vista

	// Actualizaciones
	hayActualizacion bool
	versionNueva     string
}

// NewModel crea un nuevo modelo con las dependencias inyectadas
func NewModel(container *injection.Container) (Model, error) {
	// Inputs
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

	// Menú principal
	items := createMainMenu()
	lista := createList(items, icons.Icons.App+" Shopify TUI", 0, 0)

	// Cargar tiendas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tiendas, _ := container.ListStoresUseCase.Execute(ctx)

	return Model{
		container:        container,
		vista:            VistaMenu,
		lista:            lista,
		inputNombre:      inputNombre,
		inputURL:         inputURL,
		inputGit:         inputGit,
		tiendas:          tiendas,
		cursorInput:      0,
		hayActualizacion: false,
	}, nil
}

// itemMenu representa un elemento de menú
type itemMenu struct {
	titulo string
	desc   string
	atajo  string
}

func (i itemMenu) Title() string       { return i.titulo }
func (i itemMenu) Description() string { return i.desc }
func (i itemMenu) FilterValue() string { return i.titulo }

// itemTienda representa una tienda en la lista
type itemTienda struct {
	tienda domain.Tienda
	indice int
}

func (i itemTienda) Title() string {
	// Verificar si hay servidor activo (necesitamos pasar el manager)
	return i.tienda.Nombre
}

func (i itemTienda) Description() string {
	metodo := icons.Icons.Download + " pull"
	if i.tienda.Metodo == domain.MethodGitClone {
		metodo = icons.Icons.Git + " git"
	}
	return i.tienda.URL + " [" + metodo + "]"
}

func (i itemTienda) FilterValue() string { return i.tienda.Nombre }

// Helper functions

func createMainMenu() []list.Item {
	return []list.Item{
		itemMenu{
			titulo: icons.Icons.Login + " Iniciar sesión",
			desc:   "Autenticarte en Shopify",
			atajo:  "a",
		},
		itemMenu{
			titulo: icons.Icons.Add + " Agregar tienda",
			desc:   "Registrar tienda y descargar tema",
			atajo:  "t",
		},
		itemMenu{
			titulo: icons.Icons.Server + " Desarrollo local",
			desc:   "Iniciar servidor",
			atajo:  "d",
		},
		itemMenu{
			titulo: icons.Icons.Logs + " Servidores activos",
			desc:   "Ver y administrar procesos",
			atajo:  "v",
		},
	}
}

func createStoreList(tiendas []domain.Tienda) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = itemTienda{tienda: t, indice: i + 1}
	}
	return items
}

func createDownloadMethods() []list.Item {
	return []list.Item{
		itemMenu{
			titulo: icons.Icons.Download + " Shopify Pull",
			desc:   "Desde Shopify directo",
			atajo:  "s",
		},
		itemMenu{
			titulo: icons.Icons.Git + " Git Clone",
			desc:   "Desde repositorio Git",
			atajo:  "g",
		},
	}
}

func createList(items []list.Item, titulo string, ancho, alto int) list.Model {
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

func (m *Model) recreateMainMenu() {
	items := createMainMenu()
	m.lista = createList(items, icons.Icons.App+" Shopify TUI", m.ancho, m.alto)
}
