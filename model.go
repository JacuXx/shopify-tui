// model.go - Define el estado de la aplicación
// Este archivo contiene todas las estructuras de datos que representan
// el estado actual de la app (qué pantalla está activa, qué tiendas hay, etc.)

package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Vista representa la pantalla actual de la aplicación
// En Go, los "enums" se hacen con constantes tipo iota
type Vista int

const (
	VistaMenu              Vista = iota // 0 - Menú principal
	VistaAgregarTienda                  // 1 - Formulario para agregar tienda (nombre + url)
	VistaSeleccionarMetodo              // 2 - Elegir método: Shopify Pull o Git Clone
	VistaInputGit                       // 3 - Input para URL del repositorio git
	VistaSeleccionarTienda              // 4 - Lista de tiendas para theme dev
	VistaSeleccionarModo                // 5 - Elegir modo: background o interactivo
	VistaLogs                           // 6 - Ver logs del servidor en tiempo real
	VistaServidores                     // 7 - Ver y gestionar servidores activos
)

// MetodoDescarga indica cómo se obtienen los archivos del tema
type MetodoDescarga int

const (
	MetodoShopifyPull MetodoDescarga = iota // Usar shopify theme pull
	MetodoGitClone                          // Usar git clone
)

// Tienda representa una tienda de Shopify guardada
// Los tags `json:"..."` indican cómo se guarda en el archivo JSON
type Tienda struct {
	Nombre string         `json:"nombre"`            // Nombre descriptivo (ej: "Mi Tienda Principal")
	URL    string         `json:"url"`               // URL de la tienda (ej: "mi-tienda.myshopify.com")
	Ruta   string         `json:"ruta"`              // Ruta local donde están los archivos del tema
	Metodo MetodoDescarga `json:"metodo"`            // Cómo se descargaron los archivos
	GitURL string         `json:"git_url,omitempty"` // URL del repo git (si aplica)
}

// Model contiene TODO el estado de la aplicación
// En Bubbletea, este es el corazón de tu app - todo vive aquí
type Model struct {
	// Control de navegación
	vista Vista // Pantalla actual (menú, formulario, etc.)

	// Componentes de UI (de la librería bubbles)
	lista         list.Model      // Lista del menú principal y tiendas
	tiendaParaDev Tienda          // Tienda seleccionada para theme dev (antes de elegir modo)
	inputNombre textinput.Model // Input para el nombre de la tienda
	inputURL    textinput.Model // Input para la URL de la tienda
	inputGit    textinput.Model // Input para la URL del repositorio git

	// Datos
	tiendas []Tienda // Lista de tiendas guardadas

	// Estado temporal (mientras se agrega una tienda)
	tiendaTemporal Tienda         // Tienda que se está creando
	metodoElegido  MetodoDescarga // Método elegido para descargar

	// Estado de la UI
	mensaje     string // Mensaje de estado/error para mostrar al usuario
	cursorInput int    // Qué input está activo (0=nombre, 1=url)
	ancho       int    // Ancho de la terminal
	alto        int    // Alto de la terminal
	
	// Vista de logs
	logsScroll     int  // Posición de scroll en los logs
	modoSeleccion  bool // Modo selección (permite copiar texto con mouse)
}

// itemMenu representa una opción del menú principal
// Debe implementar la interface list.Item de bubbles
type itemMenu struct {
	titulo string // Texto principal
	desc   string // Descripción debajo del título
	atajo  string // Tecla de atajo (a, s, d, etc.)
}

// Estos métodos son requeridos por la interface list.Item
func (i itemMenu) Title() string       { return i.titulo }
func (i itemMenu) Description() string { return i.desc }
func (i itemMenu) FilterValue() string { return i.titulo }

// itemTienda representa una tienda en la lista de selección
type itemTienda struct {
	tienda Tienda
	indice int // Número para selección rápida (1, 2, 3...)
}

func (i itemTienda) Title() string {
	// Mostrar indicador si tiene servidor activo
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
	
	// Mostrar URL del servidor si está activo
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

// modeloInicial crea y configura el estado inicial de la aplicación
func modeloInicial() Model {
	// === Configurar input para el nombre de la tienda ===
	inputNombre := textinput.New()
	inputNombre.Placeholder = "Mi Tienda Principal"
	inputNombre.CharLimit = 50
	inputNombre.Width = 40

	// === Configurar input para la URL de Shopify ===
	inputURL := textinput.New()
	inputURL.Placeholder = "mi-tienda.myshopify.com"
	inputURL.CharLimit = 100
	inputURL.Width = 40

	// === Configurar input para URL de Git ===
	inputGit := textinput.New()
	inputGit.Placeholder = "git@github.com:usuario/tema.git o https://..."
	inputGit.CharLimit = 200
	inputGit.Width = 50

	// === Crear las opciones del menú principal ===
	items := crearMenuPrincipal()

	// === Crear la lista del menú (sin paginación) ===
	lista := crearLista(items, Icons.App+" Shopify TUI", 0, 0)

	// === Cargar tiendas guardadas del archivo JSON ===
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

// crearMenuPrincipal retorna los items del menú principal
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

// crearListaTiendas convierte el slice de tiendas en items para la lista
func crearListaTiendas(tiendas []Tienda) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = itemTienda{tienda: t, indice: i + 1}
	}
	return items
}

// crearListaMetodos crea la lista de opciones para elegir método de descarga
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

// crearListaModos crea la lista de opciones para elegir modo de servidor
func crearListaModos(tienda Tienda, tieneServidor bool) []list.Item {
	// Opciones comunes siempre disponibles
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
		// Ya hay un servidor corriendo - mostrar opciones de servidor + comunes
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
	// No hay servidor - mostrar opción para iniciar + comunes
	items := []list.Item{
		itemMenu{
			titulo: Icons.Rocket + " Iniciar",
			desc:   "Ejecutar theme dev",
			atajo:  "i",
		},
	}
	return append(items, opcionesComunes...)
}

// crearLista crea una lista configurada sin paginación (muestra todos los items)
func crearLista(items []list.Item, titulo string, ancho, alto int) list.Model {
	// Calcular altura basada en número de items (cada item ocupa ~2 líneas + header)
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
	lista.SetShowPagination(false) // Nunca mostrar dots de paginación
	return lista
}

// recrearMenuPrincipal recrea la lista del menú principal
func (m *Model) recrearMenuPrincipal() {
	items := crearMenuPrincipal()
	m.lista = crearLista(items, Icons.App+" Shopify TUI", m.ancho, m.alto)
}

