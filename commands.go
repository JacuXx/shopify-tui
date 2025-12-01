// commands.go - Ejecución de comandos externos (Shopify CLI y Git)
// Este archivo contiene las funciones que ejecutan comandos de terminal

package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// === MENSAJES PERSONALIZADOS ===

// comandoTerminadoMsg indica que un comando terminó exitosamente
type comandoTerminadoMsg struct {
	resultado string
	tienda    *Tienda // Tienda creada (si aplica)
}

// errorMsg indica que hubo un error al ejecutar un comando
type errorMsg struct {
	err error
}

// === COMANDOS DE SHOPIFY ===

// ejecutarShopifyLogin ejecuta 'shopify auth login'
func ejecutarShopifyLogin() tea.Cmd {
	cmd := exec.Command("shopify", "auth", "login")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Sesión iniciada correctamente"}
	})
}

// ejecutarShopifyPull ejecuta 'shopify theme pull' en el directorio especificado
// Descarga el tema de la tienda al directorio local
func ejecutarShopifyPull(storeURL string, directorio string) tea.Cmd {
	// shopify theme pull --store mi-tienda.myshopify.com --path ./
	cmd := exec.Command("shopify", "theme", "pull", "--store", storeURL, "--path", directorio)
	cmd.Dir = directorio // Ejecutar en el directorio de la tienda

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Tema descargado correctamente"}
	})
}

// ejecutarThemeDev ejecuta 'shopify theme dev' en el directorio de la tienda
func ejecutarThemeDev(storeURL string, directorio string) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "dev", "--store", storeURL)
	cmd.Dir = directorio // MUY IMPORTANTE: ejecutar en el directorio del tema

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Servidor de desarrollo cerrado"}
	})
}

// === COMANDOS DE GIT ===

// ejecutarGitClone ejecuta 'git clone' para clonar un repositorio
func ejecutarGitClone(gitURL string, directorio string) tea.Cmd {
	// git clone <url> <directorio>
	// Usamos "." para clonar en el directorio actual (que ya creamos)
	cmd := exec.Command("git", "clone", gitURL, ".")
	cmd.Dir = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Repositorio clonado correctamente"}
	})
}

// === COMANDOS COMBINADOS (para el flujo de agregar tienda) ===

// ejecutarDescargaTema ejecuta la descarga del tema según el método elegido
// y retorna un mensaje con la tienda creada para agregarla a la lista
func ejecutarDescargaTema(tienda Tienda) tea.Cmd {
	return func() tea.Msg {
		// Crear el directorio para la tienda
		directorio, err := crearDirectorioTienda(tienda.Nombre)
		if err != nil {
			return errorMsg{err: err}
		}

		// Actualizar la ruta de la tienda
		tienda.Ruta = directorio

		// Crear el comando según el método
		var cmd *exec.Cmd
		if tienda.Metodo == MetodoGitClone {
			cmd = exec.Command("git", "clone", tienda.GitURL, ".")
		} else {
			cmd = exec.Command("shopify", "theme", "pull", "--store", tienda.URL, "--path", ".")
		}
		cmd.Dir = directorio

		// Ejecutar el comando
		// Nota: Aquí NO usamos ExecProcess porque queremos retornar la tienda
		// En su lugar, ejecutamos directamente y manejamos el resultado
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return errorMsg{err: err}
			}
			return comandoTerminadoMsg{
				resultado: "✅ Tienda configurada correctamente",
				tienda:    &tienda,
			}
		})()
	}
}

// ejecutarDescargaConExec ejecuta la descarga y toma control de la terminal
func ejecutarDescargaConExec(tienda Tienda, directorio string) tea.Cmd {
	var cmd *exec.Cmd
	if tienda.Metodo == MetodoGitClone {
		cmd = exec.Command("git", "clone", tienda.GitURL, ".")
	} else {
		cmd = exec.Command("shopify", "theme", "pull", "--store", tienda.URL, "--path", ".")
	}
	cmd.Dir = directorio

	// Guardar la tienda para después
	t := tienda
	t.Ruta = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado: "✅ Tienda configurada correctamente",
			tienda:    &t,
		}
	})
}

// ejecutarThemeDevInteractivo ejecuta 'shopify theme dev' en modo interactivo
// Esto toma control de la terminal para mostrar los logs
func ejecutarThemeDevInteractivo(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "dev", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: "✅ Servidor de desarrollo cerrado"}
	})
}
