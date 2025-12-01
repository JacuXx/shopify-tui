// icons.go - Sistema de iconos con Nerd Fonts y fallback ASCII
// Los iconos Nerd Font requieren una fuente compatible instalada
// Si no se detecta, usa caracteres ASCII compatibles universalmente

package main

import (
	"os"
	"strings"
)

// IconSet contiene todos los iconos usados en la aplicación
type IconSet struct {
	// Estado
	Success    string // ✓ Éxito
	Error      string // ✗ Error
	Warning    string // ⚠ Advertencia
	Info       string // ℹ Info
	
	// Acciones
	Login      string // Iniciar sesión
	Add        string // Agregar
	Delete     string // Eliminar
	Exit       string // Salir
	
	// Tienda/Desarrollo
	Store      string // Tienda
	Server     string // Servidor
	ServerOn   string // Servidor activo
	ServerOff  string // Servidor inactivo
	Logs       string // Ver logs
	Stop       string // Detener
	
	// Archivos/Git
	Download   string // Descargar
	Upload     string // Subir
	Git        string // Git
	Folder     string // Carpeta
	
	// Herramientas
	Terminal   string // Terminal
	Editor     string // Editor
	Code       string // Código
	
	// UI
	Arrow      string // Flecha
	Dot        string // Punto
	Play       string // Play/Iniciar
	Rocket     string // Cohete/Lanzar
	
	// Títulos
	App        string // Título de la app
}

// Iconos Nerd Font (requiere fuente Nerd Font instalada)
var NerdIcons = IconSet{
	// Estado
	Success:    "",  // nf-fa-check
	Error:      "",  // nf-fa-times
	Warning:    "",  // nf-fa-warning
	Info:       "",  // nf-fa-info_circle
	
	// Acciones
	Login:      "",  // nf-fa-sign_in
	Add:        "",  // nf-fa-plus
	Delete:     "",  // nf-fa-trash
	Exit:       "",  // nf-md-exit_to_app
	
	// Tienda/Desarrollo
	Store:      "",  // nf-fa-shopping_cart
	Server:     "",  // nf-fa-server
	ServerOn:   "",  // nf-fa-circle (verde)
	ServerOff:  "",  // nf-fa-circle_o
	Logs:       "",  // nf-fa-list_alt
	Stop:       "",  // nf-fa-stop
	
	// Archivos/Git
	Download:   "",  // nf-fa-download
	Upload:     "",  // nf-fa-upload
	Git:        "",  // nf-dev-git_branch
	Folder:     "",  // nf-fa-folder_open
	
	// Herramientas
	Terminal:   "",  // nf-fa-terminal
	Editor:     "",  // nf-fa-edit
	Code:       "",  // nf-fa-code
	
	// UI
	Arrow:      "",  // nf-fa-chevron_right
	Dot:        "",  // nf-fa-circle
	Play:       "",  // nf-fa-play
	Rocket:     "",  // nf-fa-rocket
	
	// Títulos
	App:        "",  // nf-fa-shopping_cart
}

// Iconos ASCII (compatibilidad universal)
var ASCIIIcons = IconSet{
	// Estado
	Success:    "[OK]",
	Error:      "[X]",
	Warning:    "[!]",
	Info:       "[i]",
	
	// Acciones
	Login:      "[>]",
	Add:        "[+]",
	Delete:     "[-]",
	Exit:       "[Q]",
	
	// Tienda/Desarrollo
	Store:      "[S]",
	Server:     "[#]",
	ServerOn:   "(*)  ",
	ServerOff:  "( )",
	Logs:       "[=]",
	Stop:       "[X]",
	
	// Archivos/Git
	Download:   "[v]",
	Upload:     "[^]",
	Git:        "[G]",
	Folder:     "[D]",
	
	// Herramientas
	Terminal:   "[$]",
	Editor:     "[E]",
	Code:       "[<>]",
	
	// UI
	Arrow:      ">",
	Dot:        "*",
	Play:       "[>]",
	Rocket:     "[!]",
	
	// Títulos
	App:        ">>",
}

// Icons es el set de iconos activo (se configura al inicio)
var Icons = ASCIIIcons

// DetectNerdFont intenta detectar si hay una Nerd Font disponible
// Verifica la variable de entorno TERM y algunos indicadores comunes
func DetectNerdFont() bool {
	// Verificar si el usuario forzó Nerd Fonts
	if env := os.Getenv("SHOPIFY_TUI_ICONS"); env != "" {
		return strings.ToLower(env) == "nerd" || strings.ToLower(env) == "true"
	}
	
	// Verificar si el usuario forzó ASCII
	if env := os.Getenv("SHOPIFY_TUI_ASCII"); env != "" {
		return false
	}
	
	// Verificar indicadores comunes de Nerd Fonts
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	
	// Terminales que comúnmente usan Nerd Fonts
	nerdTerminals := []string{
		"kitty", "alacritty", "wezterm", "iterm", "iterm2",
		"hyper", "warp", "ghostty",
	}
	
	for _, t := range nerdTerminals {
		if strings.Contains(strings.ToLower(term), t) ||
			strings.Contains(strings.ToLower(termProgram), t) {
			return true
		}
	}
	
	// Por defecto, usar Nerd Fonts en sistemas modernos
	// El usuario puede cambiar con SHOPIFY_TUI_ASCII=1
	return true
}

// InitIcons inicializa el set de iconos basado en la detección
func InitIcons() {
	if DetectNerdFont() {
		Icons = NerdIcons
	} else {
		Icons = ASCIIIcons
	}
}

// Funciones helper para construir strings con iconos
func IconSuccess(msg string) string {
	return Icons.Success + " " + msg
}

func IconError(msg string) string {
	return Icons.Error + " " + msg
}

func IconWarning(msg string) string {
	return Icons.Warning + " " + msg
}

func IconInfo(msg string) string {
	return Icons.Info + " " + msg
}
