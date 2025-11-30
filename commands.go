// commands.go - Ejecución de comandos externos (Shopify CLI)
// Este archivo contiene las funciones que ejecutan comandos de terminal
// como 'shopify auth login' y 'shopify theme dev'

package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// === MENSAJES PERSONALIZADOS ===
// En Bubbletea, los comandos retornan "mensajes" que se procesan en Update()
// Definimos tipos personalizados para cada resultado posible

// comandoTerminadoMsg indica que un comando terminó exitosamente
type comandoTerminadoMsg struct {
	resultado string // Mensaje para mostrar al usuario
}

// errorMsg indica que hubo un error al ejecutar un comando
type errorMsg struct {
	err error // El error que ocurrió
}

// === FUNCIONES DE COMANDOS ===

// ejecutarShopifyLogin ejecuta 'shopify auth login'
// Usa tea.ExecProcess para permitir interacción con el browser
// (el login de Shopify abre el navegador para OAuth)
func ejecutarShopifyLogin() tea.Cmd {
	// exec.Command crea un comando listo para ejecutar
	cmd := exec.Command("shopify", "auth", "login")

	// tea.ExecProcess es ESPECIAL: 
	// - Pausa la TUI temporalmente
	// - Deja que el comando use la terminal directamente
	// - Cuando termina, retorna el mensaje que le indiquemos
	// Esto es necesario para comandos interactivos como el login
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Sesión iniciada correctamente"}
	})
}

// ejecutarThemeDev ejecuta 'shopify theme dev --store URL'
// Este comando inicia el servidor de desarrollo local
func ejecutarThemeDev(storeURL string) tea.Cmd {
	// Crear el comando con los argumentos
	cmd := exec.Command("shopify", "theme", "dev", "--store", storeURL)

	// También usamos ExecProcess porque theme dev es interactivo:
	// - Muestra logs en tiempo real
	// - Responde a Ctrl+C para cerrar
	// - Puede pedir confirmaciones
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Servidor de desarrollo cerrado"}
	})
}

// === NOTA SOBRE tea.ExecProcess vs tea.Cmd normal ===
//
// tea.ExecProcess: Para comandos INTERACTIVOS
// - El comando toma control de la terminal
// - El usuario puede interactuar (escribir, ver output en vivo)
// - La TUI se "pausa" mientras el comando corre
// - Ideal para: login, theme dev, vim, cualquier cosa interactiva
//
// tea.Cmd normal (con exec.Command().Run()): Para comandos SILENCIOSOS  
// - El comando corre en background
// - No se ve el output
// - La TUI sigue activa
// - Ideal para: verificar versiones, operaciones rápidas
//
// En nuestro caso, TODOS los comandos de Shopify son interactivos,
// así que siempre usamos tea.ExecProcess
