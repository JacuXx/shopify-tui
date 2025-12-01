// commands.go - Ejecución de comandos externos (Shopify CLI y Git)
// Este archivo contiene las funciones que ejecutan comandos de terminal

package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// === MENSAJES PERSONALIZADOS ===

// comandoTerminadoMsg indica que un comando terminó exitosamente
type comandoTerminadoMsg struct {
	resultado       string
	tienda          *Tienda // Tienda creada (si aplica)
	volverAOpciones bool    // Si debe volver a opciones de desarrollo en vez del menú
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

// ejecutarThemePull ejecuta 'shopify theme pull' para bajar cambios del tema
func ejecutarThemePull(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "pull", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       "✅ Cambios descargados correctamente",
			volverAOpciones: true,
		}
	})
}

// ejecutarThemePush ejecuta 'shopify theme push' para subir cambios al tema
func ejecutarThemePush(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "push", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       "✅ Cambios subidos correctamente",
			volverAOpciones: true,
		}
	})
}

// ejecutarAbrirEditor abre el editor de código en el directorio del tema
// Por defecto usa VS Code (code .), pero se puede cambiar
func ejecutarAbrirEditor(tienda Tienda) tea.Cmd {
	// Intentar con VS Code primero
	cmd := exec.Command("code", ".")
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       "✅ Editor abierto",
			volverAOpciones: true,
		}
	})
}

// ejecutarAbrirTerminal abre una terminal interactiva en el directorio del tema
// Usa el shell por defecto del sistema
func ejecutarAbrirTerminal(tienda Tienda) tea.Cmd {
	// Usar zsh o bash según lo que esté disponible
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "zsh"
	}

	// Mostrar mensaje de cómo salir
	fmt.Println("\n╭─────────────────────────────────────────────────╮")
	fmt.Println("│  📂 Terminal abierta en: " + tienda.Nombre)
	fmt.Println("│  💡 Escribe 'exit' o presiona Ctrl+D para volver")
	fmt.Println("╰─────────────────────────────────────────────────╯\n")

	// Crear comando que abre shell en el directorio
	cmd := exec.Command(shell)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       "✅ Terminal cerrada",
			volverAOpciones: true,
		}
	})
}
