package icons

import (
	"fmt"
	"os"
	"strings"
)

// IconSet contiene los iconos de la aplicación
type IconSet struct {
	App       string
	Server    string
	ServerOn  string
	Download  string
	Upload    string
	Git       string
	Terminal  string
	Editor    string
	Logs      string
	Stop      string
	Rocket    string
	Login     string
	Add       string
	Success   string
	Error     string
	Warning   string
	Info      string
}

var (
	// Conjunto de iconos Unicode completo
	fullIcons = IconSet{
		App:      "🛍️",
		Server:   "🖥️",
		ServerOn: "✅",
		Download: "⬇️",
		Upload:   "⬆️",
		Git:      "🌿",
		Terminal: "⌨️",
		Editor:   "✏️",
		Logs:     "📋",
		Stop:     "⏹️",
		Rocket:   "🚀",
		Login:    "🔐",
		Add:      "➕",
		Success:  "✅",
		Error:    "❌",
		Warning:  "⚠️",
		Info:     "ℹ️",
	}

	// Conjunto de iconos ASCII fallback
	asciiIcons = IconSet{
		App:      "[APP]",
		Server:   "[SRV]",
		ServerOn: "[ON]",
		Download: "[DL]",
		Upload:   "[UP]",
		Git:      "[GIT]",
		Terminal: "[TRM]",
		Editor:   "[EDT]",
		Logs:     "[LOG]",
		Stop:     "[STOP]",
		Rocket:   "[>>]",
		Login:    "[LGN]",
		Add:      "[+]",
		Success:  "[OK]",
		Error:    "[ERR]",
		Warning:  "[!]",
		Info:     "[i]",
	}

	// Iconos actuales (se determina en init)
	Icons IconSet
)

// InitIcons inicializa los iconos basados en soporte del terminal
func InitIcons() {
	// Verificar si es Windows y si soporta Unicode
	if isWindows() && !supportsUnicode() {
		Icons = asciiIcons
		return
	}

	Icons = fullIcons
}

// isWindows verifica si el sistema operativo es Windows
func isWindows() bool {
	return os.PathSeparator == '\\'
}

// supportsUnicode verifica si la terminal soporta Unicode
func supportsUnicode() bool {
	term := os.Getenv("TERM")
	codepage := os.Getenv("LANG")

	// Si TERM está configurado a algo que soporta Unicode
	if strings.Contains(term, "256color") || strings.Contains(term, "truecolor") {
		return true
	}

	// Verificar LANG
	if strings.Contains(codepage, "UTF-8") || strings.Contains(codepage, "utf8") {
		return true
	}

	// En Windows 10+, generalmente soporta Unicode
	return true
}

// IconSuccess retorna un mensaje de éxito
func IconSuccess(msg string) string {
	return fmt.Sprintf("%s %s", Icons.Success, msg)
}

// IconError retorna un mensaje de error
func IconError(msg string) string {
	return fmt.Sprintf("%s %s", Icons.Error, msg)
}

// IconWarning retorna un mensaje de advertencia
func IconWarning(msg string) string {
	return fmt.Sprintf("%s %s", Icons.Warning, msg)
}

// IconInfo retorna un mensaje de información
func IconInfo(msg string) string {
	return fmt.Sprintf("%s %s", Icons.Info, msg)
}
