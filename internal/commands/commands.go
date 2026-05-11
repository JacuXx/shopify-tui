package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacuXx/shopify-cli/internal/domain"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
)

var shellsPermitidos = map[string]bool{
	"/bin/zsh":      true,
	"/bin/bash":     true,
	"/bin/sh":       true,
	"/usr/bin/zsh":  true,
	"/usr/bin/bash": true,
	"/usr/bin/sh":   true,
}

func esGitURLValida(u string) bool {
	for _, prefix := range []string{"https://", "http://", "git://", "git@"} {
		if strings.HasPrefix(u, prefix) {
			return true
		}
	}
	return false
}

// ComandoTerminadoMsg se emite cuando un comando externo finaliza con éxito.
type ComandoTerminadoMsg struct {
	Resultado       string
	Tienda          *domain.Tienda
	VolverAOpciones bool
}

// ErrorMsg se emite cuando un comando externo falla.
type ErrorMsg struct {
	Err error
}

// EjecutarShopifyLogin inicia el flujo OAuth de Shopify.
func EjecutarShopifyLogin() tea.Cmd {
	cmd := exec.Command("shopify", "auth", "login")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{Resultado: icons.IconSuccess("Sesión iniciada correctamente")}
	})
}

// EjecutarDescargaConExec descarga el tema según el método configurado en la tienda.
func EjecutarDescargaConExec(tienda domain.Tienda, directorio string) tea.Cmd {
	var cmd *exec.Cmd
	if tienda.Metodo == domain.MetodoGitClone {
		if !esGitURLValida(tienda.GitURL) {
			return func() tea.Msg {
				return ErrorMsg{Err: fmt.Errorf("URL de repositorio inválida")}
			}
		}
		cmd = exec.Command("git", "clone", "--", tienda.GitURL, ".")
	} else {
		cmd = exec.Command("shopify", "theme", "pull", "--store", tienda.URL, "--path", ".")
	}
	cmd.Dir = directorio

	t := tienda
	t.Ruta = directorio
	esGitClone := tienda.Metodo == domain.MetodoGitClone

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			if esGitClone && strings.Contains(err.Error(), "128") {
				return ErrorMsg{Err: fmt.Errorf("git clone falló. Verifica:\n  • URL correcta\n  • Acceso al repo (SSH key o token)\n  • Para repos privados: https://github.com/usuario/repo.git")}
			}
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{
			Resultado: icons.IconSuccess("Tienda configurada correctamente"),
			Tienda:    &t,
		}
	})
}

// EjecutarThemePull descarga los cambios del tema desde Shopify.
func EjecutarThemePull(tienda domain.Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "pull", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{
			Resultado:       icons.IconSuccess("Cambios descargados correctamente"),
			VolverAOpciones: true,
		}
	})
}

// EjecutarThemePush sube los cambios del tema a Shopify.
func EjecutarThemePush(tienda domain.Tienda) tea.Cmd {
	cmd := exec.Command("shopify", "theme", "push", "--store", tienda.URL)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{
			Resultado:       icons.IconSuccess("Cambios subidos correctamente"),
			VolverAOpciones: true,
		}
	})
}

// EjecutarAbrirEditor abre VS Code en el directorio de la tienda.
func EjecutarAbrirEditor(tienda domain.Tienda) tea.Cmd {
	cmd := exec.Command("code", ".")
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{
			Resultado:       icons.IconSuccess("Editor abierto"),
			VolverAOpciones: true,
		}
	})
}

// EjecutarAbrirTerminal abre una terminal interactiva en el directorio de la tienda.
func EjecutarAbrirTerminal(tienda domain.Tienda) tea.Cmd {
	shell := os.Getenv("SHELL")
	if !shellsPermitidos[shell] {
		shell = "/bin/sh"
	}

	fmt.Println("\n╭─────────────────────────────────────────────────╮")
	fmt.Println("│  " + icons.Icons.Folder + " Terminal abierta en: " + tienda.Nombre)
	fmt.Println("│  " + icons.Icons.Info + " Escribe 'exit' o presiona Ctrl+D para volver")
	fmt.Println("╰─────────────────────────────────────────────────╯")
	fmt.Println()

	cmd := exec.Command(shell)
	cmd.Dir = tienda.Ruta

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{
			Resultado:       icons.IconSuccess("Terminal cerrada"),
			VolverAOpciones: true,
		}
	})
}
