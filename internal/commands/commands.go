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

func esStoreURLValida(u string) bool {
	for _, c := range u {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return strings.HasSuffix(u, ".myshopify.com") && len(u) > len(".myshopify.com")
}

type ComandoTerminadoMsg struct {
	Resultado       string
	Tienda          *domain.Tienda
	VolverAOpciones bool
}

type ErrorMsg struct {
	Err error
}

func EjecutarShopifyLogin() tea.Cmd {
	cmd := exec.Command("shopify", "auth", "login")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ComandoTerminadoMsg{Resultado: icons.IconSuccess("Sesión iniciada correctamente")}
	})
}

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
		if !esStoreURLValida(tienda.URL) {
			return func() tea.Msg { return ErrorMsg{Err: fmt.Errorf("URL de tienda inválida")} }
		}
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

func EjecutarThemePull(tienda domain.Tienda) tea.Cmd {
	if !esStoreURLValida(tienda.URL) {
		return func() tea.Msg { return ErrorMsg{Err: fmt.Errorf("URL de tienda inválida")} }
	}
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

func EjecutarThemePush(tienda domain.Tienda) tea.Cmd {
	if !esStoreURLValida(tienda.URL) {
		return func() tea.Msg { return ErrorMsg{Err: fmt.Errorf("URL de tienda inválida")} }
	}
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
