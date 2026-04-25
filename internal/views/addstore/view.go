package addstore

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/JacuXx/shopify-cli/internal/ui/components"
	"github.com/JacuXx/shopify-cli/internal/ui/icons"
	"github.com/JacuXx/shopify-cli/internal/ui/styles"
)

// ViewAgregarTienda renderiza el formulario de agregar tienda (paso 1).
func ViewAgregarTienda(s State) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render("➕ Agregar Nueva Tienda"))
	b.WriteString("\n\n")

	b.WriteString(styles.Info.Render("Paso 1 de 2: Información básica"))
	b.WriteString("\n\n")

	if s.CursorInput == 0 {
		b.WriteString(styles.InputActivo.Render("> Nombre de la tienda:"))
	} else {
		b.WriteString(styles.Label.Render("  Nombre de la tienda:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + s.InputNombre.View())
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("    Un nombre para identificar la tienda"))
	b.WriteString("\n\n")

	if s.CursorInput == 1 {
		b.WriteString(styles.InputActivo.Render("> URL de Shopify:"))
	} else {
		b.WriteString(styles.Label.Render("  URL de Shopify:"))
	}
	b.WriteString("\n")
	b.WriteString("  " + s.InputURL.View() + lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Render(".myshopify.com"))
	b.WriteString("\n\n")

	if s.Mensaje != "" {
		b.WriteString(styles.Error.Render(s.Mensaje))
		b.WriteString("\n\n")
	}

	b.WriteString(styles.Ayuda.Render("tab: cambiar campo • enter: continuar • esc: cancelar"))

	return styles.Contenedor.Render(b.String())
}

// ViewSeleccionarMetodo renderiza la vista de selección de método de descarga.
func ViewSeleccionarMetodo(s State) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(icons.Icons.Download + " Método de descarga"))
	b.WriteString("\n\n")

	b.WriteString(styles.Label.Render("Tienda: "))
	b.WriteString(s.TiendaTemporal.Nombre)
	b.WriteString("\n")
	b.WriteString(styles.Label.Render("URL: "))
	b.WriteString(s.TiendaTemporal.URL)
	b.WriteString("\n\n")

	metodos := []components.ItemMenu{
		{Titulo: icons.Icons.Download + " Shopify Pull", Desc: "Desde Shopify directo", Atajo: "s"},
		{Titulo: icons.Icons.Git + " Git Clone", Desc: "Desde repositorio Git", Atajo: "g"},
	}

	b.WriteString(renderMetodosConAtajos(metodos, "Elige método"))
	b.WriteString("\n")

	if s.Mensaje != "" {
		b.WriteString(styles.Error.Render(s.Mensaje))
		b.WriteString("\n")
	}

	b.WriteString(styles.Ayuda.Render("[S]hopify Pull | [G]it Clone | l/enter | q: volver"))

	return b.String()
}

// ViewInputGit renderiza la vista de input de URL de repositorio Git.
func ViewInputGit(s State) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(icons.Icons.Git + " Clonar desde Git"))
	b.WriteString("\n\n")

	b.WriteString(styles.Label.Render("Tienda: "))
	b.WriteString(s.TiendaTemporal.Nombre)
	b.WriteString("\n\n")

	b.WriteString(styles.InputActivo.Render("> URL del repositorio:"))
	b.WriteString("\n")
	b.WriteString("  " + s.InputGit.View())
	b.WriteString("\n\n")

	b.WriteString(styles.Ayuda.Render("Ejemplos:"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  SSH:    git@github.com:usuario/tema.git"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  HTTPS:  https://github.com/usuario/tema.git"))
	b.WriteString("\n")
	b.WriteString(styles.Ayuda.Render("  Privado: https://TOKEN@github.com/usuario/tema.git"))
	b.WriteString("\n\n")

	if s.Mensaje != "" {
		b.WriteString(styles.Error.Render(s.Mensaje))
		b.WriteString("\n\n")
	}

	if s.GitURLConfirmada {
		b.WriteString(styles.Info.Render("enter: confirmar clone | cualquier tecla: cancelar"))
	} else {
		b.WriteString(styles.Ayuda.Render("enter: continuar | q: volver"))
	}

	return styles.Contenedor.Render(b.String())
}

func renderMetodosConAtajos(items []components.ItemMenu, titulo string) string {
	var b strings.Builder

	b.WriteString(styles.Titulo.Render(titulo))
	b.WriteString("\n\n")

	for _, item := range items {
		atajo := styles.Atajo.Render("[" + strings.ToUpper(item.Atajo) + "]")
		itemTitulo := styles.ItemNormal.Render(item.Titulo)
		desc := styles.Desc.Render(item.Desc)

		b.WriteString(fmt.Sprintf("  %s %s\n", atajo, itemTitulo))
		b.WriteString(fmt.Sprintf("      %s\n", desc))
	}

	return b.String()
}
