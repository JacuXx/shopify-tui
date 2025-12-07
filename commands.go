package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type comandoTerminadoMsg struct {
	resultado       string
	tienda          *Tienda
	volverAOpciones bool
}

type errorMsg struct {
	err error
}

func ejecutarShopifyLogin() tea.Cmd {
	cmd := exec.Command("shopify", "auth", "login")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: IconSuccess("Sesión iniciada correctamente")}
	})
}

func ejecutarShopifyPull(storeURL string, directorio string) tea.Cmd {

	cmd := exec.Command("shopify", "theme", "pull", "--store", storeURL, "--path", directorio)
	cmd.Dir = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: IconSuccess("Tema descargado correctamente")}
	})
}

func ejecutarThemeDev(storeURL string, directorio string) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "dev", "--store", storeURL)
	cmd.Dir = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: IconSuccess("Servidor de desarrollo cerrado")}
	})
}

func ejecutarGitClone(gitURL string, directorio string) tea.Cmd {

	cmd := exec.Command("git", "clone", gitURL, ".")
	cmd.Dir = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: IconSuccess("Repositorio clonado correctamente")}
	})
}

func ejecutarDescargaTema(tienda Tienda) tea.Cmd {
	return func() tea.Msg {

		directorio, err := crearDirectorioTienda(tienda.Nombre)
		if err != nil {
			return errorMsg{err: err}
		}

		tienda.Ruta = directorio

		var cmd *exec.Cmd
		if tienda.Metodo == MetodoGitClone {
			cmd = exec.Command("git", "clone", tienda.GitURL, ".")
		} else {
			cmd = exec.Command("shopify", "theme", "pull", "--store", tienda.URL, "--path", ".")
		}
		cmd.Dir = directorio

		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return errorMsg{err: err}
			}
			return comandoTerminadoMsg{
				resultado: IconSuccess("Tienda configurada correctamente"),
				tienda:    &tienda,
			}
		})()
	}
}

func ejecutarDescargaConExec(tienda Tienda, directorio string) tea.Cmd {
	var cmd *exec.Cmd
	if tienda.Metodo == MetodoGitClone {
		cmd = exec.Command("git", "clone", tienda.GitURL, ".")
	} else {
		cmd = exec.Command("shopify", "theme", "pull", "--store", tienda.URL, "--path", ".")
	}
	cmd.Dir = directorio

	t := tienda
	t.Ruta = directorio

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado: IconSuccess("Tienda configurada correctamente"),
			tienda:    &t,
		}
	})
}

func ejecutarThemeDevInteractivo(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "dev", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{resultado: IconSuccess("Servidor de desarrollo cerrado")}
	})
}

func ejecutarThemePull(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "pull", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       IconSuccess("Cambios descargados correctamente"),
			volverAOpciones: true,
		}
	})
}

func ejecutarThemePush(tienda Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "push", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       IconSuccess("Cambios subidos correctamente"),
			volverAOpciones: true,
		}
	})
}

func ejecutarAbrirEditor(tienda Tienda) tea.Cmd {

	cmd := exec.Command("code", ".")
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       IconSuccess("Editor abierto"),
			volverAOpciones: true,
		}
	})
}

func ejecutarAbrirTerminal(tienda Tienda) tea.Cmd {

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "zsh"
	}

	fmt.Println("\n╭─────────────────────────────────────────────────╮")
	fmt.Println("│  " + Icons.Folder + " Terminal abierta en: " + tienda.Nombre)
	fmt.Println("│  " + Icons.Info + " Escribe 'exit' o presiona Ctrl+D para volver")
	fmt.Println("╰─────────────────────────────────────────────────╯\n")

	cmd := exec.Command(shell)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return errorMsg{err: err}
		}
		return comandoTerminadoMsg{
			resultado:       IconSuccess("Terminal cerrada"),
			volverAOpciones: true,
		}
	})
}
