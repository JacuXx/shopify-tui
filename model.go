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
	VistaMenu             Vista = iota // 0 - Menú principal
	VistaAgregarTienda                 // 1 - Formulario para agregar tienda
	VistaSeleccionarTienda             // 2 - Lista de tiendas para theme dev
)

// Tienda representa una tienda de Shopify guardada
// Los tags `json:"..."` indican cómo se guarda en el archivo JSON
type Tienda struct {
	Nombre string `json:"nombre"` // Nombre descriptivo (ej: "Mi Tienda Principal")
	URL    string `json:"url"`    // URL de la tienda (ej: "mi-tienda.myshopify.com")
}

// Model contiene TODO el estado de la aplicación
// En Bubbletea, este es el corazón de tu app - todo vive aquí
type Model struct {
	// Control de navegación
	vista Vista // Pantalla actual (menú, formulario, etc.)

	// Componentes de UI (de la librería bubbles)
	lista       list.Model      // Lista del menú principal y tiendas
	inputNombre textinput.Model // Input para el nombre de la tienda
	inputURL    textinput.Model // Input para la URL de la tienda

	// Datos
	tiendas []Tienda // Lista de tiendas guardadas

	// Estado de la UI
	mensaje     string // Mensaje de estado/error para mostrar al usuario
	cursorInput int    // Qué input está activo (0=nombre, 1=url)
	ancho       int    // Ancho de la terminal
	alto        int    // Alto de la terminal
}

// itemMenu representa una opción del menú principal
// Debe implementar la interface list.Item de bubbles
type itemMenu struct {
	titulo string // Texto principal
	desc   string // Descripción debajo del título
}

// Estos métodos son requeridos por la interface list.Item
// Title() retorna el texto principal del item
func (i itemMenu) Title() string { return i.titulo }

// Description() retorna la descripción secundaria
func (i itemMenu) Description() string { return i.desc }

// FilterValue() es usado para la búsqueda/filtrado (usamos el título)
func (i itemMenu) FilterValue() string { return i.titulo }

// itemTienda representa una tienda en la lista de selección
// También implementa list.Item
type itemTienda struct {
	tienda Tienda
}

func (i itemTienda) Title() string       { return i.tienda.Nombre }
func (i itemTienda) Description() string { return i.tienda.URL }
func (i itemTienda) FilterValue() string { return i.tienda.Nombre }

// modeloInicial crea y configura el estado inicial de la aplicación
// Esta función se llama una sola vez al inicio
func modeloInicial() Model {
	// === Configurar input para el nombre de la tienda ===
	inputNombre := textinput.New()
	inputNombre.Placeholder = "Mi Tienda Principal" // Texto gris de ejemplo
	inputNombre.CharLimit = 50                      // Máximo 50 caracteres
	inputNombre.Width = 30                          // Ancho del input

	// === Configurar input para la URL ===
	inputURL := textinput.New()
	inputURL.Placeholder = "mi-tienda.myshopify.com"
	inputURL.CharLimit = 100
	inputURL.Width = 30

	// === Crear las opciones del menú principal ===
	items := []list.Item{
		itemMenu{
			titulo: "🔐 Iniciar sesión en Shopify",
			desc:   "Abre el navegador para autenticarte con tu cuenta",
		},
		itemMenu{
			titulo: "➕ Agregar tienda",
			desc:   "Guardar una nueva tienda para acceso rápido",
		},
		itemMenu{
			titulo: "🚀 Ejecutar theme dev",
			desc:   "Iniciar servidor de desarrollo local",
		},
		itemMenu{
			titulo: "❌ Salir",
			desc:   "Cerrar la aplicación (o presiona q)",
		},
	}

	// === Crear la lista del menú ===
	// list.NewDefaultDelegate() crea un renderizador por defecto
	lista := list.New(items, list.NewDefaultDelegate(), 50, 14)
	lista.Title = "🛒 Shopify TUI"
	lista.SetShowStatusBar(false)  // Ocultar barra de estado
	lista.SetFilteringEnabled(false) // Desactivar filtrado con /

	// === Cargar tiendas guardadas del archivo JSON ===
	tiendas, _ := cargarTiendas() // Si hay error, tiendas será []

	return Model{
		vista:       VistaMenu,
		lista:       lista,
		inputNombre: inputNombre,
		inputURL:    inputURL,
		tiendas:     tiendas,
		cursorInput: 0,
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

// recrearMenuPrincipal recrea la lista del menú principal
// Útil cuando volvemos al menú desde otra vista
func (m *Model) recrearMenuPrincipal() {
	items := []list.Item{
		itemMenu{
			titulo: "🔐 Iniciar sesión en Shopify",
			desc:   "Abre el navegador para autenticarte con tu cuenta",
		},
		itemMenu{
			titulo: "➕ Agregar tienda",
			desc:   "Guardar una nueva tienda para acceso rápido",
		},
		itemMenu{
			titulo: "🚀 Ejecutar theme dev",
			desc:   "Iniciar servidor de desarrollo local",
		},
		itemMenu{
			titulo: "❌ Salir",
			desc:   "Cerrar la aplicación (o presiona q)",
		},
	}

	m.lista = list.New(items, list.NewDefaultDelegate(), 50, 14)
	m.lista.Title = "🛒 Shopify TUI"
	m.lista.SetShowStatusBar(false)
	m.lista.SetFilteringEnabled(false)
}
