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
	logsScroll int // Posición de scroll en los logs
}

// itemMenu representa una opción del menú principal
// Debe implementar la interface list.Item de bubbles
type itemMenu struct {
	titulo string // Texto principal
	desc   string // Descripción debajo del título
}

// Estos métodos son requeridos por la interface list.Item
func (i itemMenu) Title() string       { return i.titulo }
func (i itemMenu) Description() string { return i.desc }
func (i itemMenu) FilterValue() string { return i.titulo }

// itemTienda representa una tienda en la lista de selección
type itemTienda struct {
	tienda Tienda
}

func (i itemTienda) Title() string {
	// Mostrar indicador si tiene servidor activo
	if ObtenerGestor().TieneServidorActivo(i.tienda.Nombre) {
		return "🟢 " + i.tienda.Nombre
	}
	return i.tienda.Nombre
}
func (i itemTienda) Description() string {
	metodo := "📥 pull"
	if i.tienda.Metodo == MetodoGitClone {
		metodo = "🔗 git"
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

	// === Crear la lista del menú ===
	lista := list.New(items, list.NewDefaultDelegate(), 50, 14)
	lista.Title = "🛒 Shopify TUI"
	lista.SetShowStatusBar(false)
	lista.SetFilteringEnabled(false)

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
			titulo: "🔐 Iniciar sesión en Shopify",
			desc:   "Abre el navegador para autenticarte con tu cuenta",
		},
		itemMenu{
			titulo: "➕ Agregar tienda",
			desc:   "Guardar una nueva tienda y descargar su tema",
		},
		itemMenu{
			titulo: "🚀 Ejecutar theme dev",
			desc:   "Iniciar servidor de desarrollo (elige modo)",
		},
		itemMenu{
			titulo: "📺 Ver servidores activos",
			desc:   "Gestionar servidores corriendo en background",
		},
		itemMenu{
			titulo: "❌ Salir",
			desc:   "Cerrar la aplicación (detiene servidores en background)",
		},
	}
}

// crearListaTiendas convierte el slice de tiendas en items para la lista
func crearListaTiendas(tiendas []Tienda) []list.Item {
	items := make([]list.Item, len(tiendas))
	for i, t := range tiendas {
		items[i] = itemTienda{tienda: t}
	}
	return items
}

// crearListaMetodos crea la lista de opciones para elegir método de descarga
func crearListaMetodos() []list.Item {
	return []list.Item{
		itemMenu{
			titulo: "📥 Shopify Pull",
			desc:   "Descargar tema directamente desde Shopify (shopify theme pull)",
		},
		itemMenu{
			titulo: "🔗 Git Clone",
			desc:   "Clonar desde un repositorio Git (SSH o HTTPS)",
		},
	}
}

// crearListaModos crea la lista de opciones para elegir modo de servidor
func crearListaModos(tienda Tienda, tieneServidor bool) []list.Item {
	if tieneServidor {
		// Ya hay un servidor corriendo - mostrar opciones para ver logs o detenerlo
		return []list.Item{
			itemMenu{
				titulo: "📺 Ver logs en vivo",
				desc:   "Ver los logs del servidor (q o esc para volver)",
			},
			itemMenu{
				titulo: "🛑 Detener servidor",
				desc:   "Terminar el servidor de desarrollo actual",
			},
		}
	}
	// No hay servidor - mostrar opciones para iniciar
	return []list.Item{
		itemMenu{
			titulo: "🔄 Iniciar en background",
			desc:   "El servidor corre en segundo plano, puedes seguir usando el CLI",
		},
		itemMenu{
			titulo: "🖥️ Iniciar interactivo",
			desc:   "Ver logs directamente en la terminal (Ctrl+C para volver al menú)",
		},
	}
}

// recrearMenuPrincipal recrea la lista del menú principal
func (m *Model) recrearMenuPrincipal() {
	items := crearMenuPrincipal()
	m.lista = list.New(items, list.NewDefaultDelegate(), 50, 16)
	m.lista.Title = "🛒 Shopify TUI"
	m.lista.SetShowStatusBar(false)
	m.lista.SetFilteringEnabled(false)
}
